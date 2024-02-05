import {
    Avatar,
    Button,
    Dropdown,
    DropdownItem,
    DropdownMenu,
    DropdownTrigger,
    Input,
    Link,
    Navbar,
    NavbarBrand,
    NavbarContent,
    NavbarMenu,
    NavbarMenuItem,
    NavbarMenuToggle
} from '@nextui-org/react';

export default function App() {
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
        <Navbar isBordered>
            <NavbarContent as="div" className="items-center" justify="end">
                <NavbarContent className="sm:hidden" justify="start">
                    <NavbarMenuToggle />
                </NavbarContent>
                <NavbarContent justify="start">
                    <NavbarBrand className="mr-2">
                        <p className="sm:block font-bold text-inherit">
                            VoxPopuli
                        </p>
                    </NavbarBrand>
                </NavbarContent>

                <Input
                    classNames={{
                        base: 'max-w-full sm:max-w-[10rem] h-10',
                        mainWrapper: 'h-full',
                        input: 'text-small',
                        inputWrapper: 'h-full font-normal text-default-500 '
                    }}
                    placeholder="Type to search..."
                    size="sm"
                    type="search"
                />
                <div className="navend">
                    <Button
                        className="navbtn"
                        as={Link}
                        color="warning"
                        href="#"
                        variant="flat"
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
                            <DropdownItem key="profile" className="h-14 gap-2">
                                <p className="font-semibold">Signed in as</p>
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
    );
}
