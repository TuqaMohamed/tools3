
FROM node:20

WORKDIR /ToolsProject/clinic-appointment-app

COPY package*.json ./

RUN npm install -g @angular/cli

RUN npm install

COPY . .

COPY replace-config-url.sh .

RUN chmod +x replace-config-url.sh

RUN ng build --configuration=production



ARG API_URL=http://localhost:8080
ENV API_URL=${API_URL}
ENV FPORT=4200

EXPOSE $FPORT


CMD ./replace-config-url.sh && ng serve --host 0.0.0.0 --configuration=production 
