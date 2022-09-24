DROP TABLE IF EXISTS client CASCADE;
DROP TABLE IF EXISTS service CASCADE;
DROP TABLE IF EXISTS appointment CASCADE;


CREATE TABLE client (
    client_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(16) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

CREATE TABLE service (
    service_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(16) NOT NULL,
    base_duration interval NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

CREATE TABLE appointment (
    client_id UUID,
    service_id UUID,
    appointment_time timestamp NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
        CONSTRAINT fk_client
        FOREIGN KEY(client_id)
        REFERENCES client(client_id),
    CONSTRAINT fk_service
        FOREIGN KEY(service_id)
        REFERENCES service(service_id)
);

INSERT INTO client (client_id, name)
VALUES
       ('1e7e0f13-50b9-41cf-9662-8f9f7abdc3b0','Мария'),
       ('db1ce0b5-6e7c-46fc-a2bb-d05b7f7db513','Владимир'),
       ('b74e2a22-2b1f-4a7b-8dfc-f08f04ba6ca2','Ольга');


INSERT INTO service (service_id, name, base_duration)
VALUES
    ('4ccfefd1-d6fe-43b9-9fe7-4eae8e198885','Массаж простаты', '00:40:00'),
    ('edba8dfe-45da-4d66-814b-6a28c3e2117c','Тайский массаж', '02:00:00'),
    ('de7df2a9-486b-4c1b-b70c-4ac3fe3e61a8','Эпиляция', '01:30:00');

INSERT INTO appointment (client_id, service_id, appointment_time)
VALUES
    ('1e7e0f13-50b9-41cf-9662-8f9f7abdc3b0', 'de7df2a9-486b-4c1b-b70c-4ac3fe3e61a8','2022-09-19 17:00:00'),
    ('db1ce0b5-6e7c-46fc-a2bb-d05b7f7db513', '4ccfefd1-d6fe-43b9-9fe7-4eae8e198885','2022-09-19 12:30:00'),
    ('b74e2a22-2b1f-4a7b-8dfc-f08f04ba6ca2', 'de7df2a9-486b-4c1b-b70c-4ac3fe3e61a8','2022-09-19 13:40:00');

SELECT client.name, service.name, appointment_time
FROM appointment
JOIN client USING(client_id)
JOIN service USING (service_id)
ORDER BY appointment_time;


select * from service;


SELECT
    service_id,
    name,
    base_duration
FROM service
WHERE service_id IN
      (SELECT service_id FROM service
       ORDER BY created_at
       LIMIT 1 OFFSET 3)