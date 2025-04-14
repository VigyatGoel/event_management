CREATE TABLE ADMIN (
    admin_id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(255),
    password VARCHAR(255),
    email VARCHAR(255)
);

CREATE TABLE ORGANISER (
    organiser_id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(20),
    password VARCHAR(255)
);

CREATE TABLE EVENT (
    event_id INT PRIMARY KEY AUTO_INCREMENT,
    organiser_id INT,
    title VARCHAR(255),
    description TEXT,
    date DATETIME,
    location VARCHAR(255),
    FOREIGN KEY (organiser_id) REFERENCES ORGANISER(organiser_id)
);

CREATE TABLE ATTENDEE (
    attendee_id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(20),
    password VARCHAR(255)
);

CREATE TABLE REGISTRATION (
    registration_id INT PRIMARY KEY AUTO_INCREMENT,
    event_id INT,
    attendee_id INT,
    registration_date DATETIME,
    status VARCHAR(50),
    FOREIGN KEY (event_id) REFERENCES EVENT(event_id),
    FOREIGN KEY (attendee_id) REFERENCES ATTENDEE(attendee_id)
);