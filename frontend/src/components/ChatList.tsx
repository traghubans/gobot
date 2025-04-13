import React from 'react';
import styled from 'styled-components';

const Container = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
`;

const Header = styled.div`
  padding: 20px;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

const Title = styled.h1`
  margin: 0;
  font-size: 20px;
  color: #1a1a1a;
`;

const NewChatButton = styled.button`
  background-color: #007AFF;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover {
    background-color: #0056b3;
  }
`;

const ChatListContainer = styled.div`
  flex: 1;
  overflow-y: auto;
`;

const ChatItem = styled.div<{ isActive: boolean }>`
  padding: 15px 20px;
  border-bottom: 1px solid #e0e0e0;
  cursor: pointer;
  background-color: ${props => props.isActive ? '#f0f0f0' : 'transparent'};
  transition: background-color 0.2s;

  &:hover {
    background-color: ${props => props.isActive ? '#f0f0f0' : '#f8f8f8'};
  }
`;

const ChatTitle = styled.div`
  font-size: 16px;
  color: #1a1a1a;
  margin-bottom: 5px;
`;

const LastMessage = styled.div`
  font-size: 14px;
  color: #666;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
`;

const TimeStamp = styled.div`
  font-size: 12px;
  color: #999;
  margin-top: 5px;
`;

interface Chat {
  id: string;
  title: string;
  lastMessage: string;
  timestamp: Date;
}

interface ChatListProps {
  chats: Chat[];
  currentChat: string | null;
  onSelectChat: (chatId: string) => void;
  onNewChat: () => void;
}

const ChatList: React.FC<ChatListProps> = ({
  chats,
  currentChat,
  onSelectChat,
  onNewChat,
}) => {
  const formatDate = (date: Date) => {
    const today = new Date();
    const messageDate = new Date(date);

    if (today.toDateString() === messageDate.toDateString()) {
      return messageDate.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    }
    return messageDate.toLocaleDateString();
  };

  return (
    <Container>
      <Header>
        <Title>Chats</Title>
        <NewChatButton onClick={onNewChat}>New Chat</NewChatButton>
      </Header>
      <ChatListContainer>
        {chats.map(chat => (
          <ChatItem
            key={chat.id}
            isActive={chat.id === currentChat}
            onClick={() => onSelectChat(chat.id)}
          >
            <ChatTitle>{chat.title}</ChatTitle>
            <LastMessage>{chat.lastMessage}</LastMessage>
            <TimeStamp>{formatDate(chat.timestamp)}</TimeStamp>
          </ChatItem>
        ))}
      </ChatListContainer>
    </Container>
  );
};

export default ChatList; 