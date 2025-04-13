import React, { useState } from 'react';
import styled from 'styled-components';
import ChatList from './components/ChatList';
import ChatWindow from './components/ChatWindow';
import InputBar from './components/InputBar';

const AppContainer = styled.div`
  display: flex;
  height: 100vh;
  background-color: #f0f0f0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
`;

const Sidebar = styled.div`
  width: 300px;
  background-color: #ffffff;
  border-right: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
`;

const MainContent = styled.div`
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: #ffffff;
`;

const LoadingIndicator = styled.div`
  display: flex;
  align-items: center;
  padding: 10px;
  color: #666;
  font-size: 14px;
  
  &::after {
    content: '';
    width: 12px;
    height: 12px;
    margin-left: 8px;
    border: 2px solid #666;
    border-top-color: transparent;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }
`;

interface Message {
  id: string;
  text: string;
  isUser: boolean;
  timestamp: Date;
}

interface Chat {
  id: string;
  title: string;
  lastMessage: string;
  timestamp: Date;
}

const App: React.FC = () => {
  const [chats, setChats] = useState<Chat[]>([]);
  const [currentChat, setCurrentChat] = useState<string | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  const handleSendMessage = async (text: string) => {
    if (!text.trim()) return;

    const newMessage: Message = {
      id: Date.now().toString(),
      text,
      isUser: true,
      timestamp: new Date(),
    };

    setMessages([...messages, newMessage]);
    setIsLoading(true);

    try {
      const response = await fetch('http://localhost:8080/query', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ query: text }),
      });

      const data = await response.json();
      
      const botMessage: Message = {
        id: (Date.now() + 1).toString(),
        text: data.answer,
        isUser: false,
        timestamp: new Date(),
      };

      setMessages(prev => [...prev, botMessage]);
    } catch (error) {
      console.error('Error sending message:', error);
      const errorMessage: Message = {
        id: (Date.now() + 1).toString(),
        text: 'Sorry, there was an error processing your request. Please try again.',
        isUser: false,
        timestamp: new Date(),
      };
      setMessages(prev => [...prev, errorMessage]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleNewChat = () => {
    const newChat: Chat = {
      id: Date.now().toString(),
      title: 'New Chat',
      lastMessage: '',
      timestamp: new Date(),
    };

    setChats([newChat, ...chats]);
    setCurrentChat(newChat.id);
    setMessages([]);
  };

  return (
    <AppContainer>
      <Sidebar>
        <ChatList
          chats={chats}
          currentChat={currentChat}
          onSelectChat={setCurrentChat}
          onNewChat={handleNewChat}
        />
      </Sidebar>
      <MainContent>
        <ChatWindow messages={messages} />
        {isLoading && <LoadingIndicator>AI is thinking...</LoadingIndicator>}
        <InputBar onSendMessage={handleSendMessage} />
      </MainContent>
    </AppContainer>
  );
};

export default App; 