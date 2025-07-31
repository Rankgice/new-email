package service

import (
	"fmt"
	"gorm.io/gorm"
	"new-email
	"new-email/internal/model"
	"new-email/internal/types"
	"runtime"
	"gorm.io/gorm"
)

// SystemStatsService 系统统计服务
	userModel    *model.UserModel
	emailModel   *model.EmailModel
	mailboxModel *model.MailboxModel
	startTime    time.Time
	startTime     time.Time
}

// NewSystemStatsService 创建系统统计服务
func NewSystemStatsService(db *gorm.DB) *SystemStatsService {
	return &SystemStatsService{
		userModel:    model.NewUserModel(db),
		emailModel:   model.NewEmailModel(db),
		mailboxModel: model.NewMailboxModel(db),
		startTime:    time.Now(),
	}
}

// GetSystemStats 获取系统统计信息
func (s *SystemStatsService) GetSystemStats() (*types.AdminSystemStatsResp, error) {
	// 获取用户统计
	userStats, err := s.getUserStats()
	if err != nil {
		return nil, err
	}

	// 获取邮件统计
	emailStats, err := s.getEmailStats()
	if err != nil {
		return nil, err
	}

	// 获取邮箱统计
	mailboxStats, err := s.getMailboxStats()
	if err != nil {
		return nil, err
	}

	// 获取系统信息
	systemInfo := s.getSystemInfo()

	return &types.AdminSystemStatsResp{
		UserStats:    *userStats,
		EmailStats:   *emailStats,
		MailboxStats: *mailboxStats,
		SystemStats:  *systemInfo,
	}, nil
}

// getUserStats 获取用户统计
func (s *SystemStatsService) getUserStats() (*types.AdminUserStats, error) {
	var stats types.AdminUserStats

	// 总用户数
	totalUsers, err := s.userModel.Count()
	if err != nil {
		return nil, err
	}
	stats.TotalUsers = totalUsers

	// 活跃用户数（最近30天登录过的用户）
	activeUsers, err := s.userModel.CountActiveUsers(30)
	if err != nil {
		return nil, err
	}
	stats.ActiveUsers = activeUsers

	// 今日新用户数
	today := time.Now().Format("2006-01-02")
	newUsers, err := s.userModel.CountNewUsers(today)
	if err != nil {
		return nil, err
	}
	stats.NewUsers = newUsers

	// 在线用户数（这里简化处理，实际应该基于session或token）
	// 假设最近5分钟有活动的用户为在线用户
	onlineUsers, err := s.userModel.CountOnlineUsers(5)
	if err != nil {
		return nil, err
	}
	stats.OnlineUsers = onlineUsers

	return &stats, nil
}

// getEmailStats 获取邮件统计
func (s *SystemStatsService) getEmailStats() (*types.AdminEmailStats, error) {
	var stats types.AdminEmailStats

	// 总邮件数
	totalEmails, err := s.emailModel.Count()
	if err != nil {
		return nil, err
	}
	stats.TotalEmails = totalEmails

	// 今日邮件数
	today := time.Now().Format("2006-01-02")
	todayEmails, err := s.emailModel.CountByDate(today)
	if err != nil {
		return nil, err
	}
	stats.TodayEmails = todayEmails

	// 发送邮件数
	sentEmails, err := s.emailModel.CountByDirection("sent")
	if err != nil {
		return nil, err
	}
	stats.SentEmails = sentEmails

	// 接收邮件数
	receivedEmails, err := s.emailModel.CountByDirection("received")
	if err != nil {
		return nil, err
	}
	stats.ReceivedEmails = receivedEmails

	return &stats, nil
}

// getMailboxStats 获取邮箱统计
func (s *SystemStatsService) getMailboxStats() (*types.AdminMailboxStats, error) {
	var stats types.AdminMailboxStats

	// 总邮箱数
	totalMailboxes, err := s.mailboxModel.Count()
	if err != nil {
		return nil, err
	}
	stats.TotalMailboxes = totalMailboxes

	// 活跃邮箱数（状态为1的邮箱）
	activeMailboxes, err := s.mailboxModel.CountByStatus(1)
	if err != nil {
		return nil, err
	}
	stats.ActiveMailboxes = activeMailboxes

	// IMAP邮箱数（这里简化处理，假设所有邮箱都支持IMAP）
	stats.ImapMailboxes = activeMailboxes

	// POP3邮箱数（这里简化处理，假设部分邮箱支持POP3）
	stats.Pop3Mailboxes = activeMailboxes / 2

	return &stats, nil
}

// getSystemInfo 获取系统信息
func (s *SystemStatsService) getSystemInfo() *types.AdminSystemInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// 计算运行时间
	uptime := time.Since(s.startTime)
	uptimeStr := formatDuration(uptime)

	// 计算内存使用率
	memUsage := float64(m.Alloc) / float64(m.Sys) * 100

	return &types.AdminSystemInfo{
		Version:   "1.0.0",
		StartTime: s.startTime,
		Uptime:    uptimeStr,
		GoVersion: runtime.Version(),
		Platform:  runtime.GOOS + "/" + runtime.GOARCH,
		CPUUsage:  s.getCPUUsage(),
		MemUsage:  memUsage,
		DiskUsage: s.getDiskUsage(),
	}
}

// getCPUUsage 获取CPU使用率（简化处理）
func (s *SystemStatsService) getCPUUsage() float64 {
	// 这里应该实现真实的CPU使用率获取
	// 简化处理，返回一个模拟值
	return 25.5
}

// getDiskUsage 获取磁盘使用率（简化处理）
func (s *SystemStatsService) getDiskUsage() float64 {
	// 这里应该实现真实的磁盘使用率获取
	// 简化处理，返回一个模拟值
	return 45.8
}

// formatDuration 格式化时间间隔
func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%d天%d小时%d分钟", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%d小时%d分钟", hours, minutes)
	} else {
		return fmt.Sprintf("%d分钟", minutes)
	}
}

// GetUserGrowthStats 获取用户增长统计
func (s *SystemStatsService) GetUserGrowthStats(days int) ([]types.DailyStats, error) {
	var stats []types.DailyStats

	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		count, err := s.userModel.CountNewUsers(date)
		if err != nil {
			return nil, err
		}

		stats = append(stats, types.DailyStats{
			Date:  date,
			Count: count,
		})
	}

	return stats, nil
}

// GetEmailGrowthStats 获取邮件增长统计
func (s *SystemStatsService) GetEmailGrowthStats(days int) ([]types.DailyStats, error) {
	var stats []types.DailyStats

	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		count, err := s.emailModel.CountByDate(date)
		if err != nil {
			return nil, err
		}

		stats = append(stats, types.DailyStats{
			Date:  date,
			Count: count,
		})
	}

	return stats, nil
}

// GetTopUsers 获取活跃用户排行
func (s *SystemStatsService) GetTopUsers(limit int) ([]types.UserActivityStats, error) {
	// 这里应该根据用户的邮件数量、登录次数等指标来排序
	// 简化处理，直接返回用户列表
	users, _, err := s.userModel.List(model.UserListParams{
		BaseListParams: model.BaseListParams{
			Page:     1,
			PageSize: limit,
		},
	})
	if err != nil {
		return nil, err
	}

	var userStats []types.UserActivityStats
	for _, user := range users {
		// 获取用户的邮件数量
		emailCount, _ := s.emailModel.CountByUserId(user.Id)

		// 获取用户的邮箱数量
		mailboxCount, _ := s.mailboxModel.CountByUserId(user.Id)

		userStats = append(userStats, types.UserActivityStats{
			UserId:       user.Id,
			Username:     user.Username,
			EmailCount:   emailCount,
			MailboxCount: mailboxCount,
			LastLoginAt:  user.LastLoginAt,
		})
	}

	return userStats, nil
}

// GetSystemHealth 获取系统健康状态
func (s *SystemStatsService) GetSystemHealth() *types.SystemHealth {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	memUsage := float64(m.Alloc) / float64(m.Sys) * 100
	cpuUsage := s.getCPUUsage()
	diskUsage := s.getDiskUsage()

	// 计算健康分数
	healthScore := 100.0
	if memUsage > 80 {
		healthScore -= 20
	}
	if cpuUsage > 80 {
		healthScore -= 20
	}
	if diskUsage > 80 {
		healthScore -= 20
	}

	status := "healthy"
	if healthScore < 60 {
		status = "critical"
	} else if healthScore < 80 {
		status = "warning"
	}

		Status:    status,
		Score:     healthScore,
		CPUUsage:  cpuUsage,
		MemUsage:  memUsage,
		DiskUsage: diskUsage,
		Uptime:    formatDuration(time.Since(s.startTime)),
		CheckedAt: time.Now(),
		CheckedAt:   time.Now(),
	}
}
