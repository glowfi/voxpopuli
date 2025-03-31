import { Loader2 } from 'lucide-react';

interface LoadingSpinnerProps {
    title?: string;
}

const LoadingSpinner = ({ title }: LoadingSpinnerProps) => {
    return (
        <div className="flex justify-center items-center min-h-fit">
            <div className="flex gap-3 items-center justify-center">
                <Loader2 className="my-4 h-8 w-8 animate-spin" />
                <p>Loading {title}...</p>
            </div>
        </div>
    );
};

export default LoadingSpinner;
