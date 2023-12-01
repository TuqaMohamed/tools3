# Use an official Node runtime as a parent image
FROM node:20

# Set the working directory in the container
WORKDIR /ToolsProject/clinic-appointment-app

# Copy package.json and package-lock.json to the working directory
COPY package*.json ./

# Install Angular CLI globally
RUN npm install -g @angular/cli

# Install the application dependencies
RUN npm install

# Copy the application files to the working directory
COPY . .

# Build the Angular application
RUN ng build --configuration=production

# Expose the port that the application will run on
EXPOSE 4200

# Define the command to run the application
CMD ["ng", "serve", "--host", "0.0.0.0"]
