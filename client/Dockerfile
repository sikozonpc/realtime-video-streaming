FROM node:12
  # if left blank app will run with dev settings
  # to build production image run:
  # $ docker build ./client --build-args app_env=production
  ENV NPM_CONFIG_LOGLEVEL warn
  ARG app_env
  ENV NODE_ENV $app_env

  RUN mkdir -p /client
  WORKDIR /client
  COPY ./ ./

  RUN npm install

  # if dev settings will use create-react start script for hot code relaoding via docker-compose shared volume
  # if production setting will build optimized static files and serve using http-server
  CMD if [ ${NODE_ENV} = production ]; \
    then \
    npm install -g http-server && \
    npm run build && \
    cd build && \
    hs -p 3000; \
    else \
    npm run start; \
    fi

  EXPOSE 3000
