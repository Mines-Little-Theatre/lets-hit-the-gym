-- Sets up the bot to follow the hours of the Mines Student Rec Center, posting at 5am (technically 4:50am but who's counting)
INSERT INTO weekdays (id, post_hour, open_hour, close_hour) VALUES
  (1, 5, 6, 23),
  (2, 5, 6, 23),
  (3, 5, 6, 23),
  (4, 5, 6, 23),
  (5, 5, 6, 20),
  (6, 5, 9, 19),
  (7, 5, 9, 19);
