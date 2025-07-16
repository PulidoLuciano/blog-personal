-- Tabla de información personal
CREATE TABLE IF NOT EXISTS personal_info (
    id INT AUTO_INCREMENT PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    bio TEXT,
    profile_image VARCHAR(255),
    email VARCHAR(255),
    github VARCHAR(255),
    linkedin VARCHAR(255)
);