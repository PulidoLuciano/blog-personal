-- Tabla de proyectos personales
CREATE TABLE IF NOT EXISTS projects (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    url_repository VARCHAR(255),
    url_website VARCHAR(255)
);