### Before your loop starts do:

  cd := countdown.Countdown{}
  
  cd.Start(len(yourArray))

### Then inside your loop do:

  cd.Count()

### Then to get the seconds remaining do:

  cd.SecondsRemaining()

### Then to get the percentage complete do:

  cd.PercentageComplete()