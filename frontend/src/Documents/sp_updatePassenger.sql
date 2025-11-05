DELIMITER $$

DROP PROCEDURE IF EXISTS `sp_updatePassenger`$$

CREATE PROCEDURE `sp_updatePassenger` (
    IN in_id INT,
    IN in_name VARCHAR(50),
    IN in_phone VARCHAR(15),
    IN in_email VARCHAR(100),
    IN in_passportNo VARCHAR(20)
)
BEGIN
    UPDATE `Passenger`
    SET `Name` = in_name, 
        `Phone` = in_phone, 
        `Email` = in_email, 
        `PassportNo` = in_passportNo
    WHERE `PassengerID` = in_id;
END$$

DELIMITER ;