-- Insert projecto
IF NOT EXISTS (SELECT 1 FROM projects) THEN
    INSERT INTO projects (title, description, url_repository, url_website, image) VALUES (
        'SportsCalendar',
        "Schedule your favorites football team's matches on google calendar with only few clicks",
        'https://github.com/PulidoLuciano/sports-calendar',
        'https://sportscalendar.vercel.app',
        '/static/assets/images/projects/sportsCalendar.webp'
    );
END IF;