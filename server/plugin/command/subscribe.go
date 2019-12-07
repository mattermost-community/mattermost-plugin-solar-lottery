// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"
)

func (c *Command) subscribe(parameters ...string) (string, error) {
	switch {
	case len(parameters) == 0:
		return fmt.Sprintf("Subscription %s created.", "storedSub.Remote.ID"), nil

	case len(parameters) == 1 && parameters[0] == "list":
		return fmt.Sprintf("Subscriptions:%s", "utils.JSONBlock(subs)"), nil

	case len(parameters) == 1 && parameters[0] == "show":
		return fmt.Sprintf("Subscription:%s", "utils.JSONBlock(storedSub)"), nil

	case len(parameters) == 1 && parameters[0] == "renew":
		return fmt.Sprintf("Subscription %s renewed until %s", "storedSub.Remote.ID", "storedSub.Remote.ExpirationDateTime"), nil

	case len(parameters) == 1 && parameters[0] == "delete":
		return fmt.Sprintf("User's subscription deleted"), nil
	}
	return "bad syntax", nil
}
