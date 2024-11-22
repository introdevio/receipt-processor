FROM scratch
COPY build/receipt-processor /
EXPOSE 8080
ENTRYPOINT ["/receipt-processor"]