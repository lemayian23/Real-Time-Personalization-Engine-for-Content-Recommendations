DELIMITER $$

DROP PROCEDURE IF EXISTS `sp_deletePassenger`$$

CREATE PROCEDURE `sp_deletePassenger` (IN in_id INT)
BEGIN
    DELETE FROM `Passenger` WHERE `PassengerID` = in_id;
    SELECT ROW_COUNT() AS 'RowsAffected';
END$$

DELIMITER ;