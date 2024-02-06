import {
    Autocomplete,
    AutocompleteItem,
    Avatar,
    Button,
    Checkbox,
    Dropdown,
    DropdownItem,
    DropdownMenu,
    DropdownTrigger,
    Input,
    Link,
    Modal,
    ModalBody,
    ModalContent,
    ModalFooter,
    ModalHeader,
    Navbar,
    NavbarBrand,
    NavbarContent,
    NavbarMenu,
    NavbarMenuItem,
    NavbarMenuToggle,
    Switch,
    useDisclosure
} from '@nextui-org/react';
import React from 'react';
import { CiMail } from 'react-icons/ci';
import { IoMdSunny } from 'react-icons/io';
import { IoMoonSharp } from 'react-icons/io5';
import { MdOutlinePassword } from 'react-icons/md';
import '../../styles/navbar.css';
import { useParams } from 'react-router-dom';
const animals = [
    {
        label: 'Cat',
        value: 'cat',
        description: 'The second most popular pet in the world'
    },
    {
        label: 'Dog',
        value: 'dog',
        description: 'The most popular pet in the world'
    },
    {
        label: 'Elephant',
        value: 'elephant',
        description: 'The largest land animal'
    },
    { label: 'Lion', value: 'lion', description: 'The king of the jungle' },
    { label: 'Tiger', value: 'tiger', description: 'The largest cat species' },
    {
        label: 'Giraffe',
        value: 'giraffe',
        description: 'The tallest land animal'
    },
    {
        label: 'Dolphin',
        value: 'dolphin',
        description: 'A widely distributed and diverse group of aquatic mammals'
    },
    {
        label: 'Penguin',
        value: 'penguin',
        description: 'A group of aquatic flightless birds'
    },
    {
        label: 'Zebra',
        value: 'zebra',
        description: 'A several species of African equids'
    },
    {
        label: 'Shark',
        value: 'shark',
        description:
            'A group of elasmobranch fish characterized by a cartilaginous skeleton'
    },
    {
        label: 'Whale',
        value: 'whale',
        description: 'Diverse group of fully aquatic placental marine mammals'
    },
    {
        label: 'Otter',
        value: 'otter',
        description: 'A carnivorous mammal in the subfamily Lutrinae'
    },
    {
        label: 'Crocodile',
        value: 'crocodile',
        description: 'A large semiaquatic reptile'
    }
];
function App({ changeTheme, theme }: any) {
    const { isOpen, onOpen, onOpenChange } = useDisclosure();
    const menuItems = [
        'Profile',
        'Dashboard',
        'Activity',
        'Analytics',
        'System',
        'Deployments',
        'My Settings',
        'Team Settings',
        'Help & Feedback',
        'Log Out'
    ];
    return (
        <>
            <Modal
                isOpen={isOpen}
                onOpenChange={onOpenChange}
                placement="top-center"
            >
                <ModalContent>
                    {(onClose) => (
                        <>
                            <ModalHeader className="flex flex-col gap-1">
                                Log in
                            </ModalHeader>
                            <ModalBody>
                                <Input
                                    autoFocus
                                    endContent={<CiMail />}
                                    label="Email"
                                    placeholder="Enter your email"
                                    variant="bordered"
                                />
                                <Input
                                    endContent={<MdOutlinePassword />}
                                    label="Password"
                                    placeholder="Enter your password"
                                    type="password"
                                    variant="bordered"
                                />
                                <div className="flex py-2 px-1 justify-between">
                                    <Checkbox
                                        classNames={{
                                            label: 'text-small'
                                        }}
                                    >
                                        Remember me
                                    </Checkbox>
                                    <Link color="primary" href="#" size="sm">
                                        Forgot password?
                                    </Link>
                                </div>
                            </ModalBody>
                            <ModalFooter>
                                <Button
                                    color="danger"
                                    variant="flat"
                                    onPress={onClose}
                                >
                                    Close
                                </Button>
                                <Button color="primary" onPress={onClose}>
                                    Sign in
                                </Button>
                            </ModalFooter>
                        </>
                    )}
                </ModalContent>
            </Modal>
            <Navbar isBordered position="sticky">
                <NavbarContent as="div" className="items-center" justify="end">
                    <NavbarContent className="sm:hidden" justify="start">
                        <NavbarMenuToggle />
                    </NavbarContent>
                    <NavbarContent justify="start">
                        <NavbarBrand className="mr-2">
                            <p className="sm:block font-bold text-inherit">
                                <Link href="/" color="foreground">
                                    VoxPopuli
                                </Link>
                            </p>
                        </NavbarBrand>
                    </NavbarContent>
                    <Autocomplete
                        defaultItems={animals}
                        label="Search"
                        size="sm"
                    >
                        {(animal) => (
                            <AutocompleteItem key={animal.value}>
                                {animal.label}
                            </AutocompleteItem>
                        )}
                    </Autocomplete>
                    {/* <Input */}
                    {/*     classNames={{ */}
                    {/*         base: 'max-w-full sm:max-w-[10rem] h-10', */}
                    {/*         mainWrapper: 'h-full', */}
                    {/*         input: 'text-small', */}
                    {/*         inputWrapper: 'h-full font-normal text-default-500 ' */}
                    {/*     }} */}
                    {/*     placeholder="Type to search..." */}
                    {/*     size="sm" */}
                    {/*     type="search" */}
                    {/* /> */}
                    <div className="navend">
                        <Switch
                            defaultSelected
                            size="sm"
                            color="success"
                            endContent={<IoMdSunny />}
                            startContent={<IoMoonSharp />}
                            onClick={() => {
                                changeTheme();
                            }}
                        >
                            {theme.charAt(0).toUpperCase() + theme.slice(1)}
                        </Switch>
                        <Button
                            className="navbtn"
                            as={Link}
                            color="warning"
                            href="#"
                            variant="flat"
                            onPress={onOpen}
                        >
                            Log in
                        </Button>
                        <Button
                            className="navbtn"
                            as={Link}
                            color="success"
                            href="#"
                            variant="flat"
                        >
                            Sign up
                        </Button>
                        <Dropdown placement="bottom-end">
                            <DropdownTrigger>
                                <Avatar
                                    isBordered
                                    as="button"
                                    className="transition-transform"
                                    color="secondary"
                                    name="Jason Hughes"
                                    size="md"
                                    src="https://i.pravatar.cc/150?u=a042581f4e29026704d"
                                />
                            </DropdownTrigger>
                            <DropdownMenu
                                aria-label="Profile Actions"
                                variant="flat"
                            >
                                <DropdownItem
                                    key="profile"
                                    className="h-14 gap-2"
                                >
                                    <p className="font-semibold">
                                        Signed in as
                                    </p>
                                    <p className="font-semibold">
                                        zoey@example.com
                                    </p>
                                </DropdownItem>
                                <DropdownItem key="settings">
                                    My Settings
                                </DropdownItem>
                                <DropdownItem key="team_settings">
                                    Team Settings
                                </DropdownItem>
                                <DropdownItem key="analytics">
                                    Analytics
                                </DropdownItem>
                                <DropdownItem key="system">System</DropdownItem>
                                <DropdownItem key="configurations">
                                    Configurations
                                </DropdownItem>
                                <DropdownItem key="help_and_feedback">
                                    Help & Feedback
                                </DropdownItem>
                                <DropdownItem key="logout" color="danger">
                                    Log Out
                                </DropdownItem>
                            </DropdownMenu>
                        </Dropdown>
                    </div>
                </NavbarContent>
                <NavbarMenu>
                    {menuItems.map((item, index) => (
                        <NavbarMenuItem key={`${item}-${index}`}>
                            <Link
                                className="w-full"
                                color={
                                    index === 2
                                        ? 'warning'
                                        : index === menuItems.length - 1
                                          ? 'danger'
                                          : 'foreground'
                                }
                                href="#"
                                size="lg"
                            >
                                {item}
                            </Link>
                        </NavbarMenuItem>
                    ))}
                </NavbarMenu>
            </Navbar>
        </>
    );
}

export default React.memo(App);
