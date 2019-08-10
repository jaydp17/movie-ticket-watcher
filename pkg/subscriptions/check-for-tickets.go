package subscriptions

//func CheckForAvailableTickets(subscriptions []Subscription) <-chan Subscription {
//	groupOfSubscriptions := groupSimilarSubscriptions(subscriptions)
//
//}

func groupSimilarSubscriptions(subscriptions []Subscription) []groupSubscriptions {
	groupedSubscriptions := make([]groupSubscriptions, 0)
	for _, sub := range subscriptions {
		foundSimilar := false
		for i := 0; i < len(groupedSubscriptions); i++ {
			// here we just compare it with the 1st subscription in the group
			// as if the 1st one passed other others should be the same
			if groupedSubscriptions[i].subscriptions[0].IsSimilar(sub) {
				groupedSubscriptions[i].subscriptions = append(groupedSubscriptions[i].subscriptions, sub)
				foundSimilar = true
				break
			}
		}
		if !foundSimilar {
			groupedElements := []Subscription{sub}
			groupedSubscriptions = append(groupedSubscriptions, groupSubscriptions{subscriptions: groupedElements})
		}
	}
	return groupedSubscriptions
}

type groupSubscriptions struct {
	subscriptions []Subscription
}
