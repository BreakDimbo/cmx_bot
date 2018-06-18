package bot

func DailyPost() {
	status := DailyAnalyze()
	botClient.Post(status)
}

func WeeklyPost() {
	status := WeeklyAnalyze()
	botClient.Post(status)
}
