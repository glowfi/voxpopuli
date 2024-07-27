import { Rocket, Flame, BadgePlus, ArrowsUpFromLine } from 'lucide-react';
import { Button } from '../ui/button';
import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectLabel,
    SelectTrigger,
    SelectValue
} from '@/components/ui/select';

const Sortposts = () => {
    return (
        <div className="flex justify-center items-center">
            <div className="hidden md:flex justify-center items-center gap-3 w-fit">
                <Button variant={'secondary'}>
                    <div className="flex gap-1">
                        <Rocket /> Best
                    </div>
                </Button>
                <Button variant={'secondary'}>
                    <div className="flex gap-1">
                        <Flame /> Hot
                    </div>
                </Button>
                <Button variant={'secondary'}>
                    <div className="flex gap-1">
                        <BadgePlus />
                        New
                    </div>
                </Button>
                <Button variant={'secondary'}>
                    <div className="flex gap-1">
                        <ArrowsUpFromLine />
                        Top
                    </div>
                </Button>
            </div>
            <div className="md:hidden">
                <Select>
                    <SelectTrigger>
                        <SelectValue placeholder="Sort By" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectGroup>
                            <SelectLabel>Choose</SelectLabel>
                            <SelectItem value="Best">
                                <div className="flex gap-1 justify-center items-center">
                                    <Rocket className="w-6 h-6" /> Best
                                </div>
                            </SelectItem>
                            <SelectItem value="Hot">
                                <div className="flex gap-1 justify-center items-center">
                                    <Flame className="w-6 h-6" /> Hot
                                </div>
                            </SelectItem>
                            <SelectItem value="New">
                                <div className="flex gap-1 justify-center items-center">
                                    <BadgePlus className="w-6 h-6" />
                                    New
                                </div>
                            </SelectItem>
                            <SelectItem value="Top">
                                <div className="flex gap-1 justify-center items-center">
                                    <ArrowsUpFromLine className="w-6 h-6" />
                                    Top
                                </div>
                            </SelectItem>
                        </SelectGroup>
                    </SelectContent>
                </Select>
            </div>
        </div>
    );
};

export default Sortposts;
