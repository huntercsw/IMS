package middleWare

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"imsPb"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	DOMAIN_NAME           = "yinuo.com"
	PV                    = new(PermissionValidator)
	TOKEN_PERIOD_VALIDITY = int64(300) // seconds
	TOKEN_ID = "1234567890"
	TICKET_OFFICE_PORT = "27892"
)

type PermissionValidator struct {
	WhiteList map[string]interface{} `json:"whiteList"` // ipAddress/hostName of remote hosts
}

func NewImsToken(domain, remoteAddr, TokenId string, createTime int64) *imsPb.ImsToken {
	return &imsPb.ImsToken{DomainName: domain, RemoteAddr: remoteAddr, TokenId: TokenId, CreateTime: createTime}
}

func (pv *PermissionValidator) PermissionValidatorInit() {
	// TODO: initialize whiteList from DB

	PV.WhiteList = map[string]interface{}{"*.*.*.*": struct{}{}}
}

func (pv *PermissionValidator) HasPermission(remoteHost string) bool {
	if _, exist := pv.WhiteList["*.*.*.*"]; exist {
		return true
	}
	_, exist := pv.WhiteList[remoteHost]
	if exist {
		return true
	} else {
		return false
	}
}

func getRemoteIP(ctx context.Context) (remoteIp string, err error) {
	peerInfo, exist := peer.FromContext(ctx)
	if ! exist {
		err = errors.New("remote client peer information dose not exist")
		return
	}

	remoteAddr := peerInfo.Addr.String()
	if remoteAddr == "" {
		err = errors.New("remote client ip dose not exist")
		return
	}

	remoteIp = strings.Split(remoteAddr, ":")[0]
	return
}

func GetTokenWithOutPermissionValida(ctx context.Context) (token *imsPb.ImsToken, err error) {
	var remoteIp string
	remoteIp, err = getRemoteIP(ctx)
	if err != nil {
		return
	}

	return NewImsToken(DOMAIN_NAME, remoteIp, TOKEN_ID, time.Now().Unix()), nil
}

func GetTokenWithPermissionValida(ctx context.Context) (token *imsPb.ImsToken, err error) {
	var remoteIp string
	remoteIp, err = getRemoteIP(ctx)
	if err != nil {
		return
	}

	if PV.HasPermission(remoteIp) {
		return NewImsToken(DOMAIN_NAME, remoteIp, TOKEN_ID, time.Now().Unix()), nil
	} else {
		err = errors.New("GetTokenWithPermissionValida: permission denied")
		return
	}
}

//
// self_defined interceptor of token authentication, to grpc service
func TokenInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	requestMetaData, exist := metadata.FromIncomingContext(ctx)
	if !exist {
		return nil, status.Errorf(401, "no meta data in request")
	}

	if !tokenValidated(requestMetaData) {
		return nil, status.Errorf(401, "token validate failed ")
	}

	return handler(ctx, req)
}

func tokenValidated(md metadata.MD) bool {
	// TODO: create token if there is not a token

	var (
		domain, remoteAddr, tokenId, createTimeStr []string
		exist                                      bool
	)
	if domain, exist = md["domain"]; !exist {				// key in meta date is lowercase(all letters is lowercase)
		fmt.Println("Domain dose not exist")
		return false
	}
	if remoteAddr, exist = md["remote-addr"]; !exist {		// key in meta date is lowercase
		fmt.Println("RemoteAddr dose not exist")
		return false
	}
	if tokenId, exist = md["token-id"]; !exist {				// key in meta date is lowercase
		fmt.Println("TokenId dose not exist")
		return false
	}
	if createTimeStr, exist = md["create-time"]; !exist {		// key in meta date is lowercase
		fmt.Println("CreateTime dose not exist")
		return false
	}

	// TODO: other validate

	if domain[0] != DOMAIN_NAME || remoteAddr[0] == "" || tokenId[0] == "" || createTimeStr[0] == "" {
		fmt.Println(fmt.Sprintf("domain[0]:%s, remote[0]:%s, tokenId[0]:%s, crearteTime[0]:%s",domain[0], remoteAddr[0], tokenId[0], createTimeStr[0]))
		return false
	}

	//
	//if createTime, err = strconv.ParseInt(createTimeStr[0], 10, 64); err != nil {
	//	fmt.Println("create time to int64 error:", err)
	//	return false
	//} else {
	//	if time.Now().Unix()-createTime >= TOKEN_PERIOD_VALIDITY {
	//		fmt.Println("token has period validity")
	//		return false
	//	}
	//}

	return true
}

//
// Grpc service StartTicketOffice
type TicketOffice struct {}

func (to *TicketOffice) GetToken(ctx context.Context, req *imsPb.GetTokenRequest) (rsp *imsPb.ImsToken, err error) {
	return GetTokenWithPermissionValida(ctx)
}

func StartTicketOffice() {
	PV = new(PermissionValidator)
	PV.PermissionValidatorInit()

	// rpc server of TicketOffice init
	listener, err := net.Listen("tcp", "localhost:" + TICKET_OFFICE_PORT)
	if err != nil {
		fmt.Printf("TicketOffice server listen on port %s error:%v \n", TICKET_OFFICE_PORT, err)
		panic("TicketOffice server start error")
	}

	s := grpc.NewServer()
	imsPb.RegisterImsTicketOfficeServer(s, &TicketOffice{})
	reflection.Register(s)

	err = s.Serve(listener)
	if err != nil {
		fmt.Println("TicketOffice grpc server error:", err)
		panic("TicketOffice grpc server error")
	}

	listener.Accept()
}

// interface WithPerRPCCredentials,  grpc client with self-defied token
type AuthenticationInToken struct{
	Token *imsPb.ImsToken
}

func (a *AuthenticationInToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"domain":     a.Token.DomainName,
		"remote-addr": a.Token.RemoteAddr,
		"token-id":    a.Token.TokenId,
		"create-time": strconv.FormatInt(a.Token.CreateTime, 10),
	}, nil
}

func (a *AuthenticationInToken) RequireTransportSecurity() bool {
	return false
}
