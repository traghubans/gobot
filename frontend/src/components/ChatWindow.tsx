import React, { useEffect, useRef } from 'react';
import styled from 'styled-components';

const Container = styled.div`
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background-color: #f0f0f0;
`;

const MessageGroup = styled.div`
  margin: 10px 0;
  display: flex;
  flex-direction: column;
`;

const Message = styled.div<{ isUser: boolean }>`
  max-width: 70%;
  padding: 10px 15px;
  border-radius: 20px;
  margin: 2px 0;
  align-self: ${props => props.isUser ? 'flex-end' : 'flex-start'};
  background-color: ${props => props.isUser ? '#007AFF' : '#e9e9eb'};
  color: ${props => props.isUser ? 'white' : '#000000'};
  font-size: 16px;
  line-height: 1.4;
  position: relative;

  &:first-child {
    border-radius: ${props => props.isUser ? '20px 20px 5px 20px' : '20px 20px 20px 5px'};
  }

  &:not(:first-child) {
    border-radius: ${props => props.isUser ? '20px 5px 5px 20px' : '5px 20px 20px 5px'};
  }

  &:last-child {
    border-radius: ${props => props.isUser ? '20px 5px 20px 20px' : '5px 20px 20px 20px'};
  }
`;

const MessageContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const ListItem = styled.div`
  display: flex;
  gap: 8px;
  align-items: flex-start;
  padding: 4px 0;
  
  &:not(:last-child) {
    margin-bottom: 4px;
  }
`;

const NumberBullet = styled.span`
  min-width: 24px;
  height: 24px;
  border-radius: 12px;
  background-color: rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 500;
`;

const ListItemContent = styled.div`
  flex: 1;
`;

const TimeStamp = styled.div<{ isUser: boolean }>`
  font-size: 12px;
  color: #8e8e93;
  margin: ${props => props.isUser ? '2px 10px 0 0' : '2px 0 0 10px'};
  align-self: ${props => props.isUser ? 'flex-end' : 'flex-start'};
`;

interface MessageType {
  id: string;
  text: string;
  isUser: boolean;
  timestamp: Date;
}

interface ChatWindowProps {
  messages: MessageType[];
}

const ChatWindow: React.FC<ChatWindowProps> = ({ messages }) => {
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (containerRef.current) {
      containerRef.current.scrollTop = containerRef.current.scrollHeight;
    }
  }, [messages]);

  const formatTime = (date: Date) => {
    return new Date(date).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };

  const parseMessageContent = (text: string) => {
    // Split text into paragraphs
    const paragraphs = text.split('\n').filter(p => p.trim());
    
    // Process each paragraph
    return paragraphs.map((paragraph, index) => {
      // Check if it's a numbered item (starts with number followed by dot or period)
      const numberedMatch = paragraph.match(/^(\d+)[.)] (.+)/);
      if (numberedMatch) {
        return {
          type: 'numbered',
          number: numberedMatch[1],
          content: numberedMatch[2].trim()
        };
      }
      
      // Check if it's a bullet point
      if (paragraph.startsWith('• ') || paragraph.startsWith('* ')) {
        return {
          type: 'bullet',
          content: paragraph.slice(2).trim()
        };
      }
      
      // Regular paragraph
      return {
        type: 'paragraph',
        content: paragraph
      };
    });
  };

  const renderMessageContent = (text: string, isUser: boolean) => {
    const content = parseMessageContent(text);
    
    return (
      <MessageContent>
        {content.map((item, index) => {
          if (item.type === 'numbered') {
            return (
              <ListItem key={index}>
                <NumberBullet>{item.number}</NumberBullet>
                <ListItemContent>{item.content}</ListItemContent>
              </ListItem>
            );
          } else if (item.type === 'bullet') {
            return (
              <ListItem key={index}>
                <NumberBullet>•</NumberBullet>
                <ListItemContent>{item.content}</ListItemContent>
              </ListItem>
            );
          } else {
            return <div key={index}>{item.content}</div>;
          }
        })}
      </MessageContent>
    );
  };

  const groupMessagesByUser = (messages: MessageType[]) => {
    const groups: MessageType[][] = [];
    let currentGroup: MessageType[] = [];

    messages.forEach((message, index) => {
      if (index === 0 || messages[index - 1].isUser !== message.isUser) {
        if (currentGroup.length > 0) {
          groups.push(currentGroup);
        }
        currentGroup = [message];
      } else {
        currentGroup.push(message);
      }
    });

    if (currentGroup.length > 0) {
      groups.push(currentGroup);
    }

    return groups;
  };

  const messageGroups = groupMessagesByUser(messages);

  return (
    <Container ref={containerRef}>
      {messageGroups.map((group, groupIndex) => (
        <MessageGroup key={groupIndex}>
          {group.map((message, messageIndex) => (
            <React.Fragment key={message.id}>
              <Message isUser={message.isUser}>
                {renderMessageContent(message.text, message.isUser)}
              </Message>
              {messageIndex === group.length - 1 && (
                <TimeStamp isUser={message.isUser}>
                  {formatTime(message.timestamp)}
                </TimeStamp>
              )}
            </React.Fragment>
          ))}
        </MessageGroup>
      ))}
    </Container>
  );
};

export default ChatWindow; 