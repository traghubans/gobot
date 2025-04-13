import 'styled-components';

declare module 'styled-components' {
  export interface DefaultTheme {
    // Define your theme properties here
    colors: {
      primary: string;
      secondary: string;
      background: string;
      text: string;
    };
    fonts: {
      body: string;
      heading: string;
    };
    spacing: {
      small: string;
      medium: string;
      large: string;
    };
  }
} 