FROM golang as build

COPY . /project

WORKDIR /project

RUN go build -o bin/post -v ./cmd/

#================================
FROM ubuntu:20.04 AS release

COPY postschema.sql /

RUN apt-get -y update && apt-get install -y locales gnupg2
RUN locale-gen en_US.UTF-8
RUN update-locale LANG=en_US.UTF-8
ENV PGVER 12
ENV DEBIAN_FRONTEND noninteractive
ENV POSTGRES_PASSWORD=admin

RUN apt-get update -y & apt-get install -y postgresql postgresql-contrib
# Run the rest of the commands as the ""postgres" user created by the postgres-$PGVER** package when it was **apt installed
USER postgres
# Create a PostgreSQL role named ""docker"" with "docker"" as the password and
# then create a database 'docker' owned by the "docker"" role.
RUN /etc/init.d/postgresql start &&\
    psql -U postgres -d postgres -a -f /postschema.sql &&\
    /etc/init.d/postgresql stop

RUN echo "listen_addresses='*'\n" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
# Expose the PostgreSQL port
EXPOSE 5432
# Add VOLUMEs to allow backup of config, logs and databases
volume ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]
# Back to the root user
USER root

COPY --from=build /project/configs /configs

COPY --from=build /project/bin /bin/

EXPOSE 8082 8086

CMD service postgresql start && post
