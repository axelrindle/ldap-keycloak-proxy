FROM alpine:3 AS prepare

COPY ./dist/* /usr/local/bin/ldap-keycloak-proxy

RUN chmod +x /usr/local/bin/ldap-keycloak-proxy


FROM gcr.io/distroless/static-debian13:nonroot

COPY --from=prepare /usr/local/bin/ldap-keycloak-proxy /usr/local/bin/ldap-keycloak-proxy

ENTRYPOINT [ "/usr/local/bin/ldap-keycloak-proxy" ]
