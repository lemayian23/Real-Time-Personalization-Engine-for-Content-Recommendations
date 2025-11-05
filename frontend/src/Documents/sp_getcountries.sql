DELIMITER $$

CREATE PROCEDURE  sp_getcountries(
	BEGIN
		SELECT * FROM `countries`
		ORDER BY `countryname`;
	
	END$$
DELIMITER;