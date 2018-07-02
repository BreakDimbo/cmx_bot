package bot

func DailyPost() {
	status, htoot := DailyAnalyze()
	botClient.PostSpoiler(status, htoot)
}

func WeeklyPost() {
	status := WeeklyAnalyze()
	botClient.Post(status)
}

func MonthlyPost() {
	status := MonthlyAnalyze()
	botClient.Post(status)
}
