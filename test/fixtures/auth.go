package fixtures

import (
	"fmt"
	"new-email/pkg/auth"
)

// AuthHelper 认证测试辅助工具
type AuthHelper struct {
	testDB *TestDB
}
docker run --hostname=52ff12b20485 --env=MINIO_ROOT_USER=minioadmin --env=MINIO_ROOT_PASSWORD=minioadmin --env=PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin --env=MINIO_ACCESS_KEY_FILE=access_key --env=MINIO_SECRET_KEY_FILE=secret_key --env=MINIO_ROOT_USER_FILE=access_key --env=MINIO_ROOT_PASSWORD_FILE=secret_key --env=MINIO_KMS_SECRET_KEY_FILE=kms_master_key --env=MINIO_UPDATE_MINISIGN_PUBKEY=RWTx5Zr1tiHQLwG9keckT0c45M3AGeHD6IvimQHpyRywVWGbP1aVSGav --env=MINIO_CONFIG_ENV_FILE=config.env --env=MC_CONFIG_DIR=/tmp/.mc --volume=/data/minio/data:/data --volume=/data/minio/config:/root/.minio --volume=/data --network=bridge --workdir=/ -p 9000:9000 --restart=no --label='architecture=x86_64' --label='build-date=2025-07-17T11:29:46' --label='com.redhat.component=ubi9-micro-container' --label='com.redhat.license_terms=https://www.redhat.com/en/about/red-hat-end-user-license-agreements#UBI' --label='description=MinIO object storage is fundamentally different. Designed for performance and the S3 API, it is 100% open-source. MinIO is ideal for large, private cloud environments with stringent security requirements and delivers mission-critical availability across a diverse range of workloads.' --label='distribution-scope=public' --label='io.buildah.version=1.41.0-dev' --label='io.k8s.description=Very small image which doesn't install the package manager.' --label='io.k8s.display-name=Red Hat Universal Base Image 9 Micro' --label='io.openshift.expose-services=' --label='maintainer=MinIO Inc <dev@min.io>' --label='name=MinIO' --label='release=RELEASE.2025-07-23T15-54-02Z' --label='summary=MinIO is a High Performance Object Storage, API compatible with Amazon S3 cloud storage service.' --label='url=https://www.redhat.com' --label='vcs-ref=52e5685a2b36133581f9e47c9e0b896cf3cbffc0' --label='vcs-type=git' --label='vendor=MinIO Inc <dev@min.io>' --label='version=RELEASE.2025-07-23T15-54-02Z' --runtime=runc -d minio/minio:latest
// NewAuthHelper 创建认证测试辅助工具
func NewAuthHelper(testDB *TestDB) *AuthHelper {
	return &AuthHelper{
		testDB: testDB,
	}
}

// GenerateUserToken 生成用户JWT Token
func (ah *AuthHelper) GenerateUserToken(userId int64) (string, error) {
	return auth.GenerateToken(userId, "user", ah.testDB.Config.JWT.Secret, ah.testDB.Config.JWT.ExpireHours)
}

// GenerateAdminToken 生成管理员JWT Token
func (ah *AuthHelper) GenerateAdminToken(adminId int64) (string, error) {
	return auth.GenerateToken(adminId, "admin", ah.testDB.Config.JWT.Secret, ah.testDB.Config.JWT.ExpireHours)
}

// GetUserAuthHeader 获取用户认证头部
func (ah *AuthHelper) GetUserAuthHeader(userId int64) map[string]string {
	token, err := ah.GenerateUserToken(userId)
	if err != nil {
		return map[string]string{}
	}
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
}

// GetAdminAuthHeader 获取管理员认证头部
func (ah *AuthHelper) GetAdminAuthHeader(adminId int64) map[string]string {
	token, err := ah.GenerateAdminToken(adminId)
	if err != nil {
		return map[string]string{}
	}
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
}

// LoginUser 用户登录并返回token
func (ah *AuthHelper) LoginUser(server *TestServer, username, password string) (string, error) {
	loginData := map[string]string{
		"username": username,
		"password": password,
	}

	w := server.POST("/api/public/user/login", loginData)
	if !server.AssertStatus(w, 200) {
		return "", fmt.Errorf("登录失败: %s", w.Body.String())
	}

	var response struct {
		Code int `json:"code"`
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}

	err := server.ParseResponse(w, &response)
	if err != nil {
		return "", err
	}

	if response.Code != 0 {
		return "", fmt.Errorf("登录失败: code=%d", response.Code)
	}

	return response.Data.Token, nil
}

// LoginAdmin 管理员登录并返回token
func (ah *AuthHelper) LoginAdmin(server *TestServer, username, password string) (string, error) {
	loginData := map[string]string{
		"username": username,
		"password": password,
	}

	w := server.POST("/api/public/admin/login", loginData)
	if !server.AssertStatus(w, 200) {
		return "", fmt.Errorf("管理员登录失败: %s", w.Body.String())
	}

	var response struct {
		Code int `json:"code"`
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}

	err := server.ParseResponse(w, &response)
	if err != nil {
		return "", err
	}

	if response.Code != 0 {
		return "", fmt.Errorf("管理员登录失败: code=%d", response.Code)
	}

	return response.Data.Token, nil
}

// CreateTestUser 创建测试用户
func (ah *AuthHelper) CreateTestUser(username, email, password string) (int64, error) {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return 0, err
	}

	user := &struct {
		Id       int64  `gorm:"primaryKey;autoIncrement"`
		Username string `gorm:"uniqueIndex;size:50;not null"`
		Email    string `gorm:"uniqueIndex;size:100;not null"`
		Password string `gorm:"size:255;not null"`
		Nickname string `gorm:"size:50"`
		Status   int    `gorm:"default:1"`
	}{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Nickname: username,
		Status:   1,
	}

	result := ah.testDB.DB.Table("users").Create(user)
	if result.Error != nil {
		return 0, result.Error
	}

	return user.Id, nil
}

// CreateTestAdmin 创建测试管理员
func (ah *AuthHelper) CreateTestAdmin(username, email, password string) (int64, error) {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return 0, err
	}

	admin := &struct {
		Id       int64  `gorm:"primaryKey;autoIncrement"`
		Username string `gorm:"uniqueIndex;size:50;not null"`
		Email    string `gorm:"uniqueIndex;size:100;not null"`
		Password string `gorm:"size:255;not null"`
		Nickname string `gorm:"size:50"`
		Role     string `gorm:"size:20;default:admin"`
		Status   int    `gorm:"default:1"`
	}{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Nickname: username,
		Role:     "admin",
		Status:   1,
	}

	result := ah.testDB.DB.Table("admins").Create(admin)
	if result.Error != nil {
		return 0, result.Error
	}

	return admin.Id, nil
}
