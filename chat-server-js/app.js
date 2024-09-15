// const express = require("express");
const http = require("http");
const cors = require("cors");
const WebSocket = require("ws");
const { SocketLogger } = require("./logs/winston");
const { NewRoom } = require("./types/Room");

const server = http.createServer();
const wss = new WebSocket.Server({ server });

const room = NewRoom();

wss.on("connection", (ws, req) => {
    // Cookie에서 User 정보 가져오기
    const cookie = req.headers.cookie;
    const [_, user] = cookie.split("=");

    room.join(ws);
    
    ws.on("message", (message) => {
        const jsonMsg = JSON.parse(message);
        jsonMsg.Name = user;
        
        room.forwardMessage(jsonMsg);
    });
    ws.on("close", () => {
        room.leave(ws);
    });
});

const PORT = process.env.PORT || 8080;
server.listen(PORT, () => {
    SocketLogger.info(`Server is running on port ${PORT}`);
});
