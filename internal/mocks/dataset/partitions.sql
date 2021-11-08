DO $$
    DECLARE
        i int4;
        sid_min int4;
        sid_max int4;
    BEGIN
        FOR i IN 1..5
            LOOP
                sid_min := (i - 1) * 5 + 1;
                sid_max := i * 5;
                EXECUTE format('CREATE TABLE service_events_p_%s ( like service_events including all )', i );
                EXECUTE format('ALTER TABLE service_events_p_%s inherit service_events', i);
                EXECUTE format('ALTER TABLE service_events_p_%s add constraint partitioning_check check ( service_id >= %s AND service_id <= %s )', i, sid_min, sid_max );
            END LOOP;
        i := 6;
        sid_min := (i - 1) * 5 + 1;
        EXECUTE format('CREATE TABLE service_events_p_%s ( like service_events including all )', i );
        EXECUTE format('ALTER TABLE service_events_p_%s inherit service_events', i);
        EXECUTE format('ALTER TABLE service_events_p_%s add constraint partitioning_check check ( service_id >= %s )', i, sid_min );
    END ;
$$;

CREATE FUNCTION partition_for_service_events() RETURNS TRIGGER AS $$
BEGIN
    EXECUTE 'INSERT INTO ' || format( 'service_events_p_%s', 1 + div(NEW.service_id - 1, 5)) || ' VALUES ($1.*)' USING NEW;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER partition_service_events BEFORE INSERT ON service_events FOR EACH ROW EXECUTE PROCEDURE partition_for_service_events();

