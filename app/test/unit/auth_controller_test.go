package unit

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"

    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/api/client/request"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/api/client/response"
    mocks "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/test/mocks/managers"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/manager"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/api/controller"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/internal/config"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/internal/runtime"
    "go.uber.org/zap"
)

func TestAuthController_RegisterAndLogin_Unit(t *testing.T) {
    // Create mock auth manager
    m := mocks.NewMockAuthManager(t)

    // Setup managers container with mock
    mgrs := &manager.Managers{AuthManager: m}

    // Minimal runtime resource (jwt config needed by controller)
    res := runtime.Resource{
        Config: config.ApplicationConfig{JwtConfig: config.JwtConfig{
            Issuer:            "test",
            SecretKey:         "secret",
            AccessExpiration:  time.Hour,
            RefreshExpiration: 24 * time.Hour,
        }},
        Logger: zap.NewNop(),
    }

    authCtrl := controller.NewAuthController(mgrs, res)

    e := echo.New()

    // Test Register success
    regReq := request.RegisterRequest{Email: "new@example.com", Password: "password123"}
    m.EXPECT().Register(mock.Anything, regReq).Return(nil)

    b, _ := json.Marshal(regReq)
    req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(b))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    if err := authCtrl.Register(c); err != nil {
        t.Fatalf("Register handler returned error: %v", err)
    }
    require.Equal(t, http.StatusOK, rec.Code)

    var genResp response.GeneralResponse[string]
    require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &genResp))
    require.Equal(t, "registered", genResp.Data)

    // Test Login success
    loginReq := request.AuthUserRequest{Email: "new@example.com", Password: "password123"}
    username := "new@example.com"
    expected := &response.AuthResponse{Username: &username, AccessToken: "a", RefreshToken: "r", ExpiresIn: 3600, TokenType: "Bearer"}
    m.EXPECT().Login(mock.Anything, loginReq).Return(expected, nil)

    lb, _ := json.Marshal(loginReq)
    req2 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(lb))
    req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec2 := httptest.NewRecorder()
    c2 := e.NewContext(req2, rec2)

    if err := authCtrl.Login(c2); err != nil {
        t.Fatalf("Login handler returned error: %v", err)
    }
    require.Equal(t, http.StatusOK, rec2.Code)

    var authResp response.GeneralResponse[response.AuthResponse]
    require.NoError(t, json.Unmarshal(rec2.Body.Bytes(), &authResp))
    require.Equal(t, "new@example.com", *authResp.Data.Username)
}




