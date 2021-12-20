from alpine

WORKDIR /brackets

COPY dist/server /brackets/
COPY migrations/* /brackets/migrations/
#COPY brackets.env.docker /brackets

EXPOSE 3000
CMD /brackets/server > /tmp/brackets.log

