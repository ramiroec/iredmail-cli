package iredmail

import (
	"fmt"
)

func (s *Server) queryMailboxAliases(mailboxEmail string) (Forwardings, error) {
	sqlQuery := `
	SELECT address, domain, forwarding, dest_domain, active, is_alias, is_forwarding, is_list 
	FROM forwardings
	WHERE forwarding = ? AND is_alias = 1
	ORDER BY domain ASC, address ASC;`

	return s.queryForwardings(sqlQuery, mailboxEmail)
}

func (s *Server) MailboxAliasAdd(alias, email string) error {
	_, domain := parseEmail(email)
	a := fmt.Sprintf("%v@%v", alias, domain)

	mailboxExists, err := s.mailboxExists(a)
	if err != nil {
		return err
	}
	if mailboxExists {
		return fmt.Errorf("An mailbox with %v already exists", a)
	}

	aliasExists, err := s.aliasExists(a)
	if err != nil {
		return err
	}
	if aliasExists {
		return fmt.Errorf("An alias with %v already exists", a)
	}

	_, err = s.DB.Exec(`
		INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_alias, active)
		VALUES ('` + a + `', '` + email + `', '` + domain + `', '` + domain + `', 1, 1)
	`)

	return err
}

func (s *Server) MailboxAliasDelete(aliasEmail string) error {
	aliasExists, err := s.mailboxAliasExists(aliasEmail)
	if err != nil {
		return err
	}
	if !aliasExists {
		return fmt.Errorf("An alias with %v doesn't exists", aliasEmail)
	}

	_, err = s.DB.Exec(`
		DELETE FROM forwardings WHERE address = '` + aliasEmail + `' AND is_alias = 1
	`)

	return err
}

func (s *Server) MailboxAliasDeleteAll(mailboxEmail string) error {
	_, err := s.DB.Exec(`
		DELETE FROM forwardings WHERE forwarding = '` + mailboxEmail + `' AND is_alias = 1
	`)

	return err
}
