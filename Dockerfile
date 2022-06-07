
FROM scratch AS final
WORKDIR /
COPY ratelimiter /
COPY etc /etc
USER MyUser
CMD [ "/ratelimiter" ]
