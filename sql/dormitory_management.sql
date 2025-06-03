-- Bahir Dar University - Bahir Dar Institute of Technology
-- Faculty of Computing
-- Advanced Database Management System Project
-- Authors: Meklit Anteneh (1507008), Tamer Getaw (1507859), Feruza Mohamed (1506085), Naom Keneni (1507336), Eyerus Tekto (1506000)

-- Create database
CREATE DATABASE bdu_student_dormitory_management_system;
USE bdu_student_dormitory_management_system;

-- Table: building
CREATE TABLE building (
    building_id INT PRIMARY KEY AUTO_INCREMENT,
    building_name VARCHAR(50),
    gender_type CHAR(6)
);

-- Table: floor
CREATE TABLE floor (
    floor_id INT PRIMARY KEY AUTO_INCREMENT,
    floor_number INT NOT NULL,
    building_id INT,
    FOREIGN KEY (building_id) REFERENCES building(building_id)
);

-- Table: room
CREATE TABLE room (
    room_id INT PRIMARY KEY AUTO_INCREMENT,
    room_number INT NOT NULL,
    capacity INT DEFAULT 3,
    amenities VARCHAR(100),
    floor_id INT,
    FOREIGN KEY (floor_id) REFERENCES floor(floor_id)
);

-- Table: department
CREATE TABLE department (
    department_id INT PRIMARY KEY AUTO_INCREMENT,
    department_name VARCHAR(40)
);

-- Table: student
CREATE TABLE student (
    student_id INT PRIMARY KEY AUTO_INCREMENT,
    first_name VARCHAR(25),
    last_name VARCHAR(25),
    gender CHAR(6),
    room_id INT,
    department_id INT,
    building_id INT,
    FOREIGN KEY (room_id) REFERENCES room(room_id),
    FOREIGN KEY (department_id) REFERENCES department(department_id),
    FOREIGN KEY (building_id) REFERENCES building(building_id)
);

-- Table: audit_log
CREATE TABLE audit_log (
    log_id INT PRIMARY KEY AUTO_INCREMENT,
    table_name VARCHAR(50),
    operation VARCHAR(50),
    record_id INT,
    change_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    changed_by VARCHAR(50),
    details VARCHAR(255)
);

-- FUNCTION 1: Get current occupants in a room
DELIMITER $$
CREATE FUNCTION get_room_occupants(p_room_id INT)
RETURNS INT
DETERMINISTIC
BEGIN
    DECLARE occupants INT;
    SELECT COUNT(*) INTO occupants
    FROM student
    WHERE room_id = p_room_id;
    RETURN occupants;
END$$
DELIMITER ;

-- FUNCTION 2: Get available beds in a room
DELIMITER $$
CREATE FUNCTION get_available_beds(p_room_id INT)
RETURNS INT
DETERMINISTIC
BEGIN
    DECLARE capacity INT;
    DECLARE occupants INT;
    SELECT capacity INTO capacity
    FROM room
    WHERE room_id = p_room_id;
    SET occupants = get_room_occupants(p_room_id);
    RETURN GREATEST(0, capacity - occupants);
END$$
DELIMITER ;

-- FUNCTION 3: Count students in a building
DELIMITER $$
CREATE FUNCTION get_building_student_count(p_building_id INT)
RETURNS INT
DETERMINISTIC
BEGIN
    DECLARE student_count INT;
    SELECT COUNT(*) INTO student_count
    FROM student
    WHERE building_id = p_building_id;
    RETURN student_count;
END$$
DELIMITER ;

-- FUNCTION 4: Validate student gender against building
DELIMITER $$
CREATE FUNCTION is_valid_gender_for_building(p_gender CHAR(6), p_building_id INT)
RETURNS BOOLEAN
DETERMINISTIC
BEGIN
    DECLARE building_gender CHAR(6);
    SELECT gender_type INTO building_gender
    FROM building
    WHERE building_id = p_building_id;
    RETURN building_gender = p_gender;
END$$
DELIMITER ;

-- TRIGGER 1: Check gender consistency
DELIMITER $$
CREATE TRIGGER check_gender_building
BEFORE INSERT ON student
FOR EACH ROW
BEGIN
    DECLARE building_gender CHAR(6);
    SELECT gender_type INTO building_gender
    FROM building
    WHERE building_id = NEW.building_id;
    IF building_gender IS NOT NULL AND building_gender != NEW.gender THEN
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Cannot assign student to a building with a different gender';
    END IF;
END$$
DELIMITER ;

-- TRIGGER 2: Check room capacity
DELIMITER $$
CREATE TRIGGER check_room_capacity
BEFORE INSERT ON student
FOR EACH ROW
BEGIN
    DECLARE current_occupants INT;
    DECLARE room_capacity INT;
    SELECT COUNT(*) INTO current_occupants
    FROM student
    WHERE room_id = NEW.room_id;
    SELECT capacity INTO room_capacity
    FROM room
    WHERE room_id = NEW.room_id;
    IF current_occupants >= room_capacity THEN
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Room capacity exceeded';
    END IF;
END$$
DELIMITER ;

-- TRIGGER 3: Log student updates
DELIMITER $$
CREATE TRIGGER audit_student_update
AFTER UPDATE ON student
FOR EACH ROW
BEGIN
    INSERT INTO audit_log (table_name, operation, record_id, changed_by, details)
    VALUES (
        'student',
        'UPDATE',
        OLD.student_id,
        USER(),
        CONCAT('Updated student: ', OLD.first_name, ' ', OLD.last_name, 
               ' to ', NEW.first_name, ' ', NEW.last_name, 
               ', room_id: ', OLD.room_id, ' to ', NEW.room_id)
    );
END$$
DELIMITER ;

-- STORED PROCEDURE 1: Assign student to room
DELIMITER $$
CREATE PROCEDURE assign_student_to_room(
    IN p_first_name VARCHAR(25),
    IN p_last_name VARCHAR(25),
    IN p_gender CHAR(6),
    IN p_room_id INT,
    IN p_department_id INT,
    IN p_building_id INT
)
BEGIN
    IF is_valid_gender_for_building(p_gender, p_building_id) THEN
        IF get_available_beds(p_room_id) > 0 THEN
            INSERT INTO student (first_name, last_name, gender, room_id, department_id, building_id)
            VALUES (p_first_name, p_last_name, p_gender, p_room_id, p_department_id, p_building_id);
            SELECT 'Student assigned successfully' AS message;
        ELSE
            SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'No available beds in the specified room';
        END IF;
    ELSE
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Gender does not match building gender';
    END IF;
END$$
DELIMITER ;

-- STORED PROCEDURE 2: Generate occupancy report
DELIMITER $$
CREATE PROCEDURE generate_occupancy_report()
BEGIN
    SELECT 
        b.building_name,
        f.floor_number,
        r.room_number,
        r.capacity,
        get_room_occupants(r.room_id) AS current_occupants,
        get_available_beds(r.room_id) AS available_beds
    FROM building b
    JOIN floor f ON b.building_id = f.building_id
    JOIN room r ON f.floor_id = r.floor_id
    ORDER BY b.building_name, f.floor_number, r.room_number;
END$$
DELIMITER ;

-- STORED PROCEDURE 3: Update room amenities
DELIMITER $$
CREATE PROCEDURE update_room_amenities(
    IN p_room_id INT,
    IN p_amenities VARCHAR(100)
)
BEGIN
    UPDATE room
    SET amenities = p_amenities
    WHERE room_id = p_room_id;
    SELECT 'Room amenities updated successfully' AS message;
END$$
DELIMITER ;

-- SECURITY: Create roles and assign permissions
CREATE ROLE 'admin', 'manager', 'viewer';
GRANT ALL ON bdu_student_dormitory_management_system.* TO 'admin';
GRANT SELECT, INSERT, UPDATE ON bdu_student_dormitory_management_system.* TO 'manager';
GRANT SELECT ON bdu_student_dormitory_management_system.* TO 'viewer';

CREATE USER 'admin_user'@'localhost' IDENTIFIED BY 'admin_pass';
CREATE USER 'manager_user'@'localhost' IDENTIFIED BY 'manager_pass';
CREATE USER 'viewer_user'@'localhost' IDENTIFIED BY 'viewer_pass';

GRANT 'admin' TO 'admin_user'@'localhost';
GRANT 'manager' TO 'manager_user'@'localhost';
GRANT 'viewer' TO 'viewer_user'@'localhost';

SET DEFAULT ROLE 'admin' FOR 'admin_user'@'localhost';
SET DEFAULT ROLE 'manager' FOR 'manager_user'@'localhost';
SET DEFAULT ROLE 'viewer' FOR 'viewer_user'@'localhost';

-- Sample data
INSERT INTO building (building_name, gender_type)
VALUES ('Aklilu Lema', 'male'), ('Abdisa Aga', 'male'), ('Lucy', 'female'), ('Xayitu', 'female');

INSERT INTO floor (floor_number, building_id)
VALUES (1, 1), (2, 1), (1, 2), (2, 2), (1, 3), (2, 3), (1, 4), (2, 4);

INSERT INTO room (room_number, floor_id, amenities)
VALUES (101, 1, 'Wi-Fi, TV'), (102, 1, 'Wi-Fi'), (103, 1, 'Air Conditioning'), (104, 1, NULL),
       (201, 2, NULL), (202, 2, NULL), (203, 2, NULL), (204, 2, NULL),
       (101, 5, 'Wi-Fi, TV'), (102, 5, NULL), (103, 5, NULL), (104, 5, NULL),
       (201, 6, NULL), (202, 6, NULL), (203, 6, NULL), (404, 6, NULL);

INSERT INTO department (department_name)
VALUES ('Software Engineering'), ('Computer Science'), ('Information System');

INSERT INTO student (first_name, last_name, gender, room_id, department_id, building_id)
VALUES ('Nahom', 'Keneni', 'male', 1, 1, 1),
       ('Meklit', 'Abebaw', 'female', 9, 2, 3),
       ('Mikiyas', 'Arage', 'male', 2, 1, 1),
       ('Bethelhem', 'Melese', 'female', 10, 3, 3);

-- Views
CREATE VIEW students_by_building AS
SELECT b.building_name, s.first_name, s.last_name, s.gender
FROM student s
JOIN building b ON s.building_id = b.building_id;

CREATE VIEW room_occupancy AS
SELECT r.room_number, get_room_occupants(r.room_id) AS current_occupants, r.capacity
FROM room r;

CREATE VIEW department_enrollment AS
SELECT d.department_name, COUNT(s.student_id) AS student_count
FROM department d
LEFT JOIN student s ON d.department_id = s.department_id
GROUP BY d.department_name;

-- Indexes
CREATE INDEX idx_student_room ON student(room_id);
CREATE INDEX idx_floor_building ON floor(building_id);
CREATE INDEX idx_department_name ON department(department_name);

-- Sample operations
INSERT INTO building (building_name, gender_type) VALUES ('New Building', 'female');
UPDATE student SET department_id = 1 WHERE student_id = 4;
CALL update_room_amenities(1, 'Wi-Fi, TV, Air Conditioning');
CALL assign_student_to_room('Abebe', 'Kebede', 'male', 3, 1, 1);
CALL generate_occupancy_report();