package auth

import (
	"github.com/pkg/errors"
	"go_im/im/api/router"
	"go_im/im/client"
	"go_im/im/dao"
	"go_im/im/message"
)

var ResponseHandleFunc func(uid int64, device int64, message *message.Message)

func respond(uid int64, seq int64, action message.Action, data interface{}) {
	resp := message.NewMessage(seq, action, data)
	respondMessage(uid, resp)
}

func respondMessage(uid int64, msg *message.Message) {
	ResponseHandleFunc(uid, 0, msg)
}

type Interface interface {
	AuthToken(info *route.Context, req *AuthTokenReq) error
	SignIn(info *route.Context, req interface{}) error
	Logout(info *route.Context) error
	Register(info *route.Context, req interface{}) error
}

type AuthApi struct {
}

func (*AuthApi) AuthToken(info *route.Context, req *AuthTokenReq) error {
	panic("implement me")
}

func (*AuthApi) SignIn(info *route.Context, req interface{}) error {
	panic("implement me")
}

func (*AuthApi) Register(ctx *route.Context, registerEntity *RegisterRequest) error {

	resp := message.NewMessage(ctx.Seq, "", "success")
	err := dao.UserDao.AddUser(registerEntity.Account, registerEntity.Password)

	if err != nil {
		resp = message.NewMessage(ctx.Seq, "", err)
	}
	ctx.Response(resp)
	return err
}

func (a *AuthApi) Logout(info *route.Context, r *LogoutRequest) error {
	err := dao.UserDao.Logout(info.Uid, r.Device, r.Token)
	if err != nil {
		return err
	}
	client.Manager.ClientLogout(info.Uid, r.Device)
	return nil
}

func (a *AuthApi) Auth(ctx *route.Context, request *AuthRequest) error {

	var resp = message.NewMessage(ctx.Seq, "", "success")
	uid := dao.UserDao.GetUid(request.Token)
	if uid > 0 {
		client.Manager.ClientSignIn(ctx.Uid, uid, request.DeviceId)
		respondMessage(uid, resp)
		return nil
	} else {
		return errors.New("token expired")
	}
}

func (a *AuthApi) Login(msg *route.Context, request *LoginRequest) error {

	if len(request.Account) == 0 || len(request.Password) == 0 {
		return errors.New("account or password empty")
	}

	uid, token, err := dao.UserDao.GetUidByLogin(request.Account, request.Password, request.Device)
	if err != nil {
		return err
	}

	m := message.NewMessage(msg.Seq, "", "success")
	if err = m.SetData(AuthorResponse{Token: token, Uid: uid}); err != nil {
		return err
	}
	client.Manager.ClientSignIn(msg.Uid, uid, request.Device)
	msg.Response(m)
	return nil
}