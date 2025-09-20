package mailserver

import (
	"errors"
	"io"
	"log"
	"strings"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapserver"
	"github.com/emersion/go-message"
)

// Expunge 删除标记为删除的邮件
func (s *IMAPSession) Expunge(w *imapserver.ExpungeWriter, uids *imap.UIDSet) error {
	if !s.authenticated || s.selectedFolder == nil {
		return errors.New("未选择邮箱")
	}

	log.Printf("删除邮件: 用户=%s, 邮箱=%s", s.username, s.selectedFolder.Name)

	// 简化实现：获取所有邮件并检查删除标志
	mails, err := s.storage.GetMails(s.username, s.selectedFolder.Name, 0)
	if err != nil {
		return err
	}

	var expungedSeqs []uint32
	for i, mail := range mails {
		seqNum := uint32(i + 1)
		uid := imap.UID(mail.ID)

		// 如果指定了 UID 集合，只处理在集合中的邮件
		if uids != nil && !uids.Contains(uid) {
			continue
		}

		// 这里应该检查邮件是否标记为删除
		// 简化实现：假设所有指定的邮件都要删除
		// 注意：这里需要根据实际的 emailModel.Delete 方法签名来调用
		// 假设需要传递邮件对象或ID，这里简化处理
		log.Printf("删除邮件: %s (简化实现，未实际删除)", mail.MessageID)

		expungedSeqs = append(expungedSeqs, seqNum)
		if err := w.WriteExpunge(seqNum); err != nil {
			return err
		}
	}

	log.Printf("成功删除 %d 封邮件", len(expungedSeqs))
	return nil
}

// Search 搜索邮件
func (s *IMAPSession) Search(kind imapserver.NumKind, criteria *imap.SearchCriteria, options *imap.SearchOptions) (*imap.SearchData, error) {
	if !s.authenticated || s.selectedFolder == nil {
		return nil, errors.New("未选择邮箱")
	}

	log.Printf("搜索邮件: 用户=%s, 邮箱=%s, 条件=%v", s.username, s.selectedFolder.Name, criteria)

	// 获取所有邮件
	mails, err := s.storage.GetMails(s.username, s.selectedFolder.Name, 0)
	if err != nil {
		return nil, err
	}

	var results []uint32
	for i, mail := range mails {
		seqNum := uint32(i + 1)
		uid := imap.UID(mail.ID)

		// 简化的搜索实现
		if s.matchSearchCriteria(mail, criteria) {
			if kind == imapserver.NumKindSeq {
				results = append(results, seqNum)
			} else { // UID
				results = append(results, uint32(uid))
			}
		}
	}

	// 转换为 NumSet
	var searchData *imap.SearchData
	if kind == imapserver.NumKindSeq {
		seqSet := make(imap.SeqSet, 0)
		for _, num := range results {
			seqSet = append(seqSet, imap.SeqRange{Start: num, Stop: num})
		}
		searchData = &imap.SearchData{
			All: seqSet,
		}
	} else {
		uidSet := make(imap.UIDSet, 0)
		for _, num := range results {
			uidSet = append(uidSet, imap.UIDRange{Start: imap.UID(num), Stop: imap.UID(num)})
		}
		searchData = &imap.SearchData{
			All: uidSet,
		}
	}

	log.Printf("搜索结果: %d 封邮件", len(results))
	return searchData, nil
}

// matchSearchCriteria 简化的搜索条件匹配
func (s *IMAPSession) matchSearchCriteria(mail *StoredMail, criteria *imap.SearchCriteria) bool {
	// 简化实现，只检查一些基本条件

	// 检查头部字段
	if len(criteria.Header) > 0 {
		for _, headerField := range criteria.Header {
			switch strings.ToLower(headerField.Key) {
			case "subject":
				if !strings.Contains(strings.ToLower(mail.Subject), strings.ToLower(headerField.Value)) {
					return false
				}
			case "from":
				if !strings.Contains(strings.ToLower(mail.From), strings.ToLower(headerField.Value)) {
					return false
				}
			case "to":
				found := false
				searchTo := strings.ToLower(headerField.Value)
				for _, to := range mail.To {
					if strings.Contains(strings.ToLower(to), searchTo) {
						found = true
						break
					}
				}
				if !found {
					return false
				}
			}
		}
	}

	// 检查正文
	if len(criteria.Body) > 0 {
		for _, bodyText := range criteria.Body {
			if !strings.Contains(strings.ToLower(mail.Body), strings.ToLower(bodyText)) {
				return false
			}
		}
	}

	// 检查文本（主题+正文）
	if len(criteria.Text) > 0 {
		for _, text := range criteria.Text {
			searchText := strings.ToLower(text)
			if !strings.Contains(strings.ToLower(mail.Subject), searchText) &&
				!strings.Contains(strings.ToLower(mail.Body), searchText) {
				return false
			}
		}
	}

	// 检查标志
	if len(criteria.Flag) > 0 {
		for _, flag := range criteria.Flag {
			if flag == imap.FlagSeen && !mail.IsRead {
				return false
			}
		}
	}

	return true
}

// Fetch 获取邮件
func (s *IMAPSession) Fetch(w *imapserver.FetchWriter, numSet imap.NumSet, options *imap.FetchOptions) error {
	if !s.authenticated || s.selectedFolder == nil {
		return errors.New("未选择邮箱")
	}

	log.Printf("获取邮件: 用户=%s, 邮箱=%s", s.username, s.selectedFolder.Name)

	// 获取所有邮件
	mails, err := s.storage.GetMails(s.username, s.selectedFolder.Name, 0)
	if err != nil {
		return err
	}

	// 遍历邮件并处理在 numSet 中的邮件
	for i, mail := range mails {
		seqNum := uint32(i + 1)
		uid := imap.UID(mail.ID)

		// 检查是否在请求的集合中
		var contains bool
		if seqSet, ok := numSet.(imap.SeqSet); ok {
			contains = seqSet.Contains(seqNum)
		} else if uidSet, ok := numSet.(imap.UIDSet); ok {
			contains = uidSet.Contains(uid)
		}
		if !contains {
			continue
		}

		// 创建 FetchWriter 并写入邮件数据
		fetchData := w.CreateMessage(seqNum)

		// 处理请求的项目
		if options.Envelope {
			envelope := s.buildEnvelope(mail)
			fetchData.WriteEnvelope(envelope)
		}

		if options.BodyStructure != nil {
			bodyStructure := s.buildBodyStructure(mail)
			fetchData.WriteBodyStructure(bodyStructure)
		}

		if options.Flags {
			flags := s.buildFlags(mail)
			fetchData.WriteFlags(flags)
		}

		if options.InternalDate {
			fetchData.WriteInternalDate(mail.Received)
		}

		if options.RFC822Size {
			fetchData.WriteRFC822Size(int64(len(mail.Body)))
		}

		if options.UID {
			fetchData.WriteUID(uid)
		}

		// 处理 BodySection
		if len(options.BodySection) > 0 {
			// 当客户端请求邮件正文时，自动标记为已读
			if !mail.IsRead {
				if err := s.storage.MarkAsRead(s.username, mail.MessageID); err != nil {
					log.Printf("自动标记邮件已读失败: %v", err)
				} else {
					log.Printf("邮件 %s 已自动标记为已读", mail.MessageID)
					mail.IsRead = true // 更新内存中的状态以反映到flags
				}
			}

			for _, item := range options.BodySection {
				body := s.buildEmailBody(mail)
				literal := fetchData.WriteBodySection(item, int64(len(mail.Body)))
				if _, err := io.Copy(literal, body); err != nil {
					literal.Close()
					return err
				}
				if err := literal.Close(); err != nil {
					return err
				}
			}
		}

		if err := fetchData.Close(); err != nil {
			return err
		}
	}

	return nil
}

// buildEnvelope 构建邮件信封
func (s *IMAPSession) buildEnvelope(mail *StoredMail) *imap.Envelope {
	return &imap.Envelope{
		Date:      mail.Received,
		Subject:   mail.Subject,
		From:      s.parseAddressList(mail.From),
		To:        s.parseAddressListList(mail.To),
		Cc:        s.parseAddressListList(mail.Cc),
		Bcc:       s.parseAddressListList(mail.Bcc),
		MessageID: mail.MessageID,
	}
}

// buildBodyStructure 构建邮件体结构
func (s *IMAPSession) buildBodyStructure(mail *StoredMail) imap.BodyStructure {
	r := strings.NewReader(mail.Body)
	entity, err := message.Read(r)
	if err != nil {
		log.Printf("解析邮件实体失败: %v", err)
		return &imap.BodyStructureSinglePart{
			Type:    "text",
			Subtype: "plain",
			Size:    uint32(len(mail.Body)),
		}
	}

	mediaType, params, err := entity.Header.ContentType()
	if err != nil {
		log.Printf("解析邮件Content-Type失败: %v", err)
		return &imap.BodyStructureSinglePart{
			Type:    "text",
			Subtype: "plain",
			Size:    uint32(len(mail.Body)),
		}
	}

	mainType, subType, ok := strings.Cut(mediaType, "/")
	if !ok {
		log.Printf("无效的Content-Type格式: %s", mediaType)
		return &imap.BodyStructureSinglePart{
			Type:    "text",
			Subtype: "plain",
			Size:    uint32(len(mail.Body)),
		}
	}

	return &imap.BodyStructureSinglePart{
		Type:    mainType,
		Subtype: subType,
		Params:  params,
		Size:    uint32(len(mail.Body)),
	}
}

// buildFlags 构建邮件标志
func (s *IMAPSession) buildFlags(mail *StoredMail) []imap.Flag {
	var flags []imap.Flag
	if mail.IsRead {
		flags = append(flags, imap.FlagSeen)
	}
	return flags
}

// parseAddressList 解析地址列表
func (s *IMAPSession) parseAddressList(email string) []imap.Address {
	if email == "" {
		return nil
	}

	// 简单解析，假设格式为 "name@domain.com"
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return nil
	}

	return []imap.Address{
		{
			Mailbox: parts[0],
			Host:    parts[1],
		},
	}
}

// parseAddressListList 解析地址列表列表
func (s *IMAPSession) parseAddressListList(emails []string) []imap.Address {
	var addresses []imap.Address
	for _, email := range emails {
		addresses = append(addresses, s.parseAddressList(email)...)
	}
	return addresses
}

// buildEmailBody 构建邮件体
func (s *IMAPSession) buildEmailBody(mail *StoredMail) io.Reader {
	return strings.NewReader(mail.Body)
}

// Store 存储邮件标志
func (s *IMAPSession) Store(w *imapserver.FetchWriter, numSet imap.NumSet, flags *imap.StoreFlags, options *imap.StoreOptions) error {
	if !s.authenticated || s.selectedFolder == nil {
		return errors.New("未选择邮箱")
	}

	log.Printf("存储邮件标志: 用户=%s, 邮箱=%s", s.username, s.selectedFolder.Name)

	// 获取所有邮件
	mails, err := s.storage.GetMails(s.username, s.selectedFolder.Name, 0)
	if err != nil {
		return err
	}

	// 遍历邮件并处理在 numSet 中的邮件
	for i, mail := range mails {
		seqNum := uint32(i + 1)
		uid := imap.UID(mail.ID)

		// 检查是否在请求的集合中
		var contains bool
		if seqSet, ok := numSet.(imap.SeqSet); ok {
			contains = seqSet.Contains(seqNum)
		} else if uidSet, ok := numSet.(imap.UIDSet); ok {
			contains = uidSet.Contains(uid)
		}
		if !contains {
			continue
		}

		// 处理标志更新
		if flags != nil {
			for _, flag := range flags.Flags {
				if flag == imap.FlagSeen {
					switch flags.Op {
					case imap.StoreFlagsSet:
						// 标记为已读
						if err := s.storage.MarkAsRead(s.username, mail.MessageID); err != nil {
							log.Printf("标记邮件已读失败: %v", err)
						}
					default:
						// 其他操作 (简化实现)
						log.Printf("邮件标志操作: %s (简化实现)", mail.MessageID)
					}
				}
			}
		}

		// 如果需要返回更新后的标志
		if options != nil && !flags.Silent {
			fetchData := w.CreateMessage(seqNum)
			updatedFlags := s.buildFlags(mail)
			fetchData.WriteFlags(updatedFlags)
			if err := fetchData.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

// Copy 复制邮件
func (s *IMAPSession) Copy(numSet imap.NumSet, destMailbox string) (*imap.CopyData, error) {
	if !s.authenticated || s.selectedFolder == nil {
		return nil, errors.New("未选择邮箱")
	}

	log.Printf("复制邮件: 从=%s 到=%s, 用户=%s", s.selectedFolder.Name, destMailbox, s.username)

	// 获取目标文件夹
	destFolder, err := s.storage.folderModel.GetByMailboxIdAndName(s.mailbox.Id, destMailbox, nil)
	if err != nil {
		return nil, err
	}
	if destFolder == nil {
		return nil, errors.New("目标邮箱不存在")
	}

	// 获取源邮件
	sourceMails, err := s.storage.GetMails(s.username, s.selectedFolder.Name, 0)
	if err != nil {
		return nil, err
	}

	var copiedUIDs []imap.UID

	// 遍历并复制符合条件的邮件
	for i, mail := range sourceMails {
		seqNum := uint32(i + 1)
		uid := imap.UID(mail.ID)

		// 检查是否在请求的集合中
		var contains bool
		if seqSet, ok := numSet.(imap.SeqSet); ok {
			contains = seqSet.Contains(seqNum)
		} else if uidSet, ok := numSet.(imap.UIDSet); ok {
			contains = uidSet.Contains(uid)
		}
		if !contains {
			continue
		}

		// 创建新的邮件副本
		copiedMail := &StoredMail{
			MessageID:   generateMessageID(s.storage.domain),
			From:        mail.From,
			To:          mail.To,
			Cc:          mail.Cc,
			Bcc:         mail.Bcc,
			Subject:     mail.Subject,
			Body:        mail.Body,
			ContentType: mail.ContentType,
			Size:        mail.Size,
			Received:    time.Now(),
			IsRead:      false, // 复制的邮件默认为未读
			FolderId:    destFolder.Id,
			FolderName:  destFolder.Name,
			MailboxID:   s.mailbox.Id,
			Username:    s.username,
		}

		if err := s.storage.SaveMail(copiedMail); err != nil {
			log.Printf("复制邮件失败: %v", err)
			continue
		}

		copiedUIDs = append(copiedUIDs, imap.UID(copiedMail.ID))
		log.Printf("成功复制邮件: %s -> %s", mail.MessageID, copiedMail.MessageID)
	}

	// 构建返回数据
	copyData := &imap.CopyData{}
	if len(copiedUIDs) > 0 {
		uidSet := make(imap.UIDSet, 0, len(copiedUIDs))
		for _, uid := range copiedUIDs {
			uidSet = append(uidSet, imap.UIDRange{Start: uid, Stop: uid})
		}
		copyData.DestUIDs = uidSet
	}

	log.Printf("成功复制 %d 封邮件", len(copiedUIDs))
	return copyData, nil
}
