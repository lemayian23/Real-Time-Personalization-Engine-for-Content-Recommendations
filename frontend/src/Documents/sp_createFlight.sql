DELIMITER $$
CREATE PROCEDURE `sp_createFlight` (
    IN in_flightNumber VARCHAR(10),
    IN in_departureTime DATETIME,
    IN in_arrivalTime DATETIME,
    IN in_originAirportID INT,
    IN in_destinationAirportID INT,
    IN in_aircraftID INT,
    IN in_ecoRoute BOOLEAN
)
BEGIN
    INSERT INTO `Flight` (
        `FlightNumber`, 
        `DepartureTime`, 
        `ArrivalTime`, 
        `OriginAirportID`, 
        `DestinationAirportID`, 
        `AircraftID`, 
        `EcoRoute`
    ) VALUES (
        in_flightNumber,
        in_departureTime,
        in_arrivalTime,
        in_originAirportID,
        in_destinationAirportID,
        in_aircraftID,
        IFNULL(in_ecoRoute, FALSE)
    );
END$$
DELIMITER ;