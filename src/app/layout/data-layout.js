"use client";
import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import {
    BarChartOutlined,
    UploadOutlined,
    UserOutlined,
    AlignCenterOutlined
} from '@ant-design/icons';
import { Layout, Menu, theme } from 'antd';
import ConfigTable from "@/app/component/configData/configTable";
import RecordTable from "@/app/component/recordData/recordTable";
import BanIPTable from "@/app/component/banIPData/banIPTable";
import RuleTable from "@/app/component/ruleData/ruleTable";

const menuConfig = [
    { label: '数据管理', href: '/dadawre' },
    { label: '操作指南', href: '/asd' },
    // 后两项保持默认配置
];

const topMenuItems = new Array(4).fill(null).map((_, index) => ({
    key: index + 1,
    label: menuConfig[index]?.label || `nav ${index + 1}`,
    href: menuConfig[index]?.href || '/second'
}));

const { Header, Content, Footer, Sider } = Layout;

const labels = ['配置数据', '规则数据', '配置封禁IP', '日志数据'];
const siderItems = [
    UserOutlined,
    AlignCenterOutlined,
    UploadOutlined,
    BarChartOutlined,
].map((icon, index) => ({
    key: String(index + 1),
    icon: React.createElement(icon),
    label: labels[index] || `nav ${index + 1}`, // 使用逻辑或作为兜底
}));

export default function MyLayout() {
    const [collapsed, setCollapsed] = useState(false);
    const [activeKey, setActiveKey] = useState(0); // 当前选中的菜单项，默认是第一个

    const {
        token: { colorBgContainer, borderRadiusLG },
    } = theme.useToken();

    // 在 activeKey 变化时，触发重新加载数据
    useEffect(() => {
        // 此处可以放置每个表格组件的数据加载逻辑
        console.log(`切换到菜单 ${activeKey}，重新加载数据`);
    }, [activeKey]);

    const renderContent = () => {
        switch (activeKey) {
            case 1:
                return <ConfigTable key="configTable" />;
            case 2:
                return <RuleTable key="ruleTable" />;
            case 3:
                return <BanIPTable key="banIPTable" />;
            case 4:
                return <RecordTable key="recordTable" />;
            default:
                return <ConfigTable />;
        }
    };

    return (
        <Layout style={{ minHeight: '100vh' }}>
            <Sider collapsible collapsed={collapsed} onCollapse={(value) => setCollapsed(value)}>
                {/* 插入 Logo 图片 */}
                <div className="demo-logo-vertical" style={{ textAlign: 'center', padding: '22.5px 0' }}>
                    <img
                        src="/fore-end/public/next.svg"
                        alt="Logo"
                        style={{ width: '80%', maxHeight: '50px', objectFit: 'contain' }}
                    />
                </div>
                <Menu
                    theme="dark"
                    mode="inline"
                    selectedKeys={[String(activeKey)]} // 根据 activeKey 设置高亮项
                    onClick={({ key }) => setActiveKey(Number(key))} // 更新选中的菜单项
                    items={siderItems}
                />
            </Sider>
            <Layout>
                <Header style={{ display: 'flex', alignItems: 'center', padding: 0, background: colorBgContainer }}>
                    <div className="demo-logo" />
                    <Menu
                        theme="dark"
                        mode="horizontal"
                        defaultSelectedKeys={['2']}
                        style={{ flex: 1, minWidth: 0 }}
                    >
                        {topMenuItems.map((item) => (
                            <Menu.Item key={item.key}>
                                <Link href={item.href}>{item.label}</Link>
                            </Menu.Item>
                        ))}
                    </Menu>
                </Header>
                <Content style={{ margin: '24px 16px 0', overflow: 'initial' }}>
                    {renderContent()} {/* 渲染根据 activeKey 决定的组件 */}
                </Content>
                <Footer style={{ textAlign: 'center' }}>
                    Ant Design ©{new Date().getFullYear()} Created by Ant UED
                </Footer>
            </Layout>
        </Layout>
    );
}
