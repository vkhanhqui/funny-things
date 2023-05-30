# Use the official Maven 3.2 image with OpenJDK 8 as the build stage
FROM maven:3.8.4-openjdk-8 AS build

# Set the working directory
WORKDIR /app

# Copy the Maven project files to the container
COPY pom.xml .

# Copy the source code
COPY src ./src

# Build the Maven project
RUN mvn clean install

# Use the official Tomcat 8.5 image as the runtime stage
FROM tomcat:8.5.57

# Copy the built WAR file to the Tomcat webapps directory
COPY --from=build app/target/chat-web-app.war /usr/local/tomcat/webapps/

# Add scripts to run on containers
COPY app_wait.sh .
