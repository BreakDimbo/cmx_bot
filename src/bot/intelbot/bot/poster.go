package bot

func DailyPost() {
	status, pic := DailyAnalyze()
	botClient.PostWithPicture(status, pic)
}

func WeeklyPost() {
	status, pic := WeeklyAnalyze()
	botClient.PostWithPicture(status, pic)
}

func MonthlyPost() {
	status := MonthlyAnalyze()
	botClient.Post(status)
}
