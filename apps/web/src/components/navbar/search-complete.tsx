'use client';

import { useState } from 'react';
import { Input } from '@/components/ui/input';
import Link from 'next/link';
import Image from 'next/image';

export default function Component() {
    const [searchTerm, setSearchTerm] = useState('');
    const [suggestions, setSuggestions] = useState([]);
    const products = [
        {
            id: 1,
            name: 'Cozy Fleece Hoodie',
            image: '/placeholder.svg'
        },
        {
            id: 2,
            name: 'Stylish Leather Backpack',
            image: '/placeholder.svg'
        },
        {
            id: 3,
            name: 'Ergonomic Desk Chair',
            image: '/placeholder.svg'
        },
        {
            id: 4,
            name: 'High-Performance Running Shoes',
            image: '/placeholder.svg'
        },
        {
            id: 5,
            name: 'Sleek Stainless Steel Water Bottle',
            image: '/placeholder.svg'
        },
        {
            id: 6,
            name: 'Plush Memory Foam Mattress',
            image: '/placeholder.svg'
        },
        {
            id: 7,
            name: 'Versatile Outdoor Camping Gear',
            image: '/placeholder.svg'
        },
        {
            id: 8,
            name: 'Luxury Silk Bedding Set',
            image: '/placeholder.svg'
        }
    ];
    //@ts-ignore
    const handleSearch = (e) => {
        const term = e.target.value;
        setSearchTerm(term);
        const matchingProducts = products.filter((product) =>
            product.name.toLowerCase().includes(term.toLowerCase())
        );
        //@ts-ignore
        setSuggestions(matchingProducts);
    };
    //@ts-ignore
    const handleSuggestionClick = (product) => {
        window.location.href = `/product/${product.id}`;
    };
    return (
        <div className="relative flex-1">
            <Input
                type="text"
                placeholder="Start typing to search..."
                value={searchTerm}
                onChange={handleSearch}
                className="w-full rounded-full bg-muted pl-10 pr-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary text-center"
            />
            {suggestions.length > 0 && (
                <div className="absolute z-10 mt-2 w-full rounded-lg border  shadow-lg">
                    {suggestions.map((product) => (
                        <Link
                            //@ts-ignore
                            key={product.id}
                            href="#"
                            className="flex items-center gap-4 px-4 py-2 hover:bg-gray-100 dark:hover:bg-black hover:cursor-pointer dark:bg-muted"
                            onClick={() => handleSuggestionClick(product)}
                            prefetch={false}
                        >
                            <Image
                                src="/placeholder.svg"
                                //@ts-ignore
                                alt={product.name}
                                width={40}
                                height={40}
                                className="rounded-md"
                            />
                            <div className="text-sm font-medium">
                                {/* @ts-ignore */}
                                {product.name}
                            </div>
                        </Link>
                    ))}
                </div>
            )}
        </div>
    );
}
