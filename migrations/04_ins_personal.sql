-- Insert de datos personales inicial
IF NOT EXISTS (SELECT 1 FROM personal_info) THEN
    INSERT INTO personal_info (full_name, bio, profile_image, email, github, linkedin) VALUES (
        'Luciano Pulido',
    'Soy una persona proactiva, apasionada por la tecnología y el desarrollo, que aprendí de manera autodidacta y profundicé mis conocimientos a través de mis estudios universitarios. Actualmente, estoy cursando el cuarto año de ingeniería en sistemas de información, con todas las materias previas aprobadas. Disfruto enfrentar desafíos y aprender algo nuevo cada día, manteniéndome siempre en constante crecimiento y ávido de nuevas experiencias.',
        'assets/images/me.webp',
        'example@example.com',
        'PulidoLuciano',
        'lucianopulido'
    );
END IF;