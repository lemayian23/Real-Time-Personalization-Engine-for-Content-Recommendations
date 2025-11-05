DELIMITER $$

DROP PROCEDURE IF EXISTS `sp_createPassenger`$$

CREATE PROCEDURE `sp_createPassenger` (
    IN in_name VARCHAR(50),
    IN in_phone VARCHAR(15),
    IN in_email VARCHAR(100),
    IN in_passportNo VARCHAR(20)
)
BEGIN
    INSERT INTO `Passenger` (`Name`, `Phone`, `Email`, `PassportNo`)
    VALUES (in_name, in_phone, in_email, in_passportNo);
END$$

DELIMITER ;