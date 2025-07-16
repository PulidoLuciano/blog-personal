-- Insert de datos personales inicial
IF NOT EXISTS (SELECT 1 FROM personal_info) THEN
    INSERT INTO personal_info (full_name, description, bio, profile_image, x, github, linkedin) VALUES (
        'Luciano Pulido',
        'Estudiante y desarrollador',
        'Soy Luciano Nicolás Pulido. Nací y vivo en Tucumán, Argentina. Desde niño me apasionaron las computadoras. Aprendí mucho por mi cuenta durante la secundaria y, al graduarme, elegí Ingeniería de Sistemas. Actualmente estudio en la Universidad Tecnológica Nacional - Facultad Regional Tucumán. Estoy desarrollando He practicado con diversas tecnologías y lenguajes. Siempre estoy aprendiendo nuevas habilidades. Me encanta leer y jugar videojuegos. ¡Estoy deseando conocer tus propuestas!',
        '/static/assets/images/me.webp',
        'https://x.com/luciano_pulido',
        'https://github.com/PulidoLuciano',
        'https://www.linkedin.com/in/lucianopulido/'
    );
END IF;