DELIMITER $$

DROP PROCEDURE IF EXISTS `sp_readPassenger`$$

CREATE PROCEDURE `sp_readPassenger` (IN in_id INT)
BEGIN
    SELECT * FROM `Passenger` WHERE `PassengerID` = in_id;
END$$

DELIMITER ;