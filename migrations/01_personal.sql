-- Tabla de información personal
CREATE TABLE IF NOT EXISTS personal_info (
    id INT AUTO_INCREMENT PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    bio TEXT NOT NULL,
    profile_image VARCHAR(255) NOT NULL,
    x VARCHAR(255) NOT NULL,
    github VARCHAR(255) NOT NULL,
    linkedin VARCHAR(255) NOT NULL
);