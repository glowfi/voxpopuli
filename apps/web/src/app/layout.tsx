import type { Metadata } from 'next';
import './globals.css';
import { Provider } from '@voxpopuli/trpc-client/src/Provider';
import Navbar from '@/components/navbar/navbar';
import { ThemeProvider } from 'next-themes';

export const metadata: Metadata = {
    title: 'VoxPopuli',
    description: 'A social media app'
};

export default function RootLayout({
    children
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="en">
            <body>
                <Provider>
                    <ThemeProvider
                        attribute="class"
                        defaultTheme="system"
                        enableSystem
                        disableTransitionOnChange
                    >
                        <Navbar />
                        {children}
                    </ThemeProvider>
                </Provider>
            </body>
        </html>
    );
}
