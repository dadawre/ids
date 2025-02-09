"use client";
import React, { useEffect, useRef, useState } from 'react';
import { Button, Input, Modal, Table, Form, Popconfirm, Space } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import Highlighter from 'react-highlight-words';

export default function BanIPTable(){
    const [data, setData] = useState([]);
    const [total, setTotal] = useState(0); // 总数据量
    const [searchText, setSearchText] = useState('');
    const [searchedColumn, setSearchedColumn] = useState('');
    const [pagination, setPagination] = useState({
        current: 1,
        pageSize: 10,
    });
    const [openModal, setOpenModal] = useState(false); // 控制模态框打开
    const [form] = Form.useForm();

    // 获取数据并更新状态
    const fetchData = async (query = '', page = 1, pageSize = 10) => {
        const response = await fetch(`http://localhost:8080/banIPData?query=${query}&page=${page}&pageSize=${pageSize}`);
        const result = await response.json();
        setData(result.data);
        setTotal(result.total); // 设置总数据量
    };

    useEffect(() => {
        fetchData(searchText, pagination.current, pagination.pageSize); // 初次加载数据
    }, [pagination.current, pagination.pageSize, searchText]);

    // 处理搜索
    const handleSearch = (selectedKeys, confirm, dataIndex) => {
        confirm();
        setSearchText(selectedKeys[0]);
        setSearchedColumn(dataIndex);
        fetchData(selectedKeys[0], pagination.current, pagination.pageSize); // 按搜索条件重新获取数据
    };

    const handleReset = (clearFilters) => {
        clearFilters();
        setSearchText('');
        fetchData('', pagination.current, pagination.pageSize); // 清除搜索条件时重新获取所有数据
    };

    // 获取列搜索功能配置
    const getColumnSearchProps = (dataIndex) => ({
        filterDropdown: ({ setSelectedKeys, selectedKeys, confirm, clearFilters }) => (
            <div
                style={{
                    padding: 8,
                }}
                onKeyDown={(e) => e.stopPropagation()}
            >
                <Input
                    placeholder={`Search ${dataIndex}`}
                    value={selectedKeys[0]}
                    onChange={(e) => setSelectedKeys(e.target.value ? [e.target.value] : [])}
                    onPressEnter={() => handleSearch(selectedKeys, confirm, dataIndex)}
                    style={{
                        marginBottom: 8,
                        display: 'block',
                    }}
                />
                <Space>
                    <Button
                        type="primary"
                        onClick={() => handleSearch(selectedKeys, confirm, dataIndex)}
                        icon={<SearchOutlined />}
                        size="small"
                        style={{
                            width: 90,
                        }}
                    >
                        Search
                    </Button>
                    <Button
                        onClick={() => clearFilters && handleReset(clearFilters)}
                        size="small"
                        style={{
                            width: 90,
                        }}
                    >
                        Reset
                    </Button>
                </Space>
            </div>
        ),
        filterIcon: (filtered) => (
            <SearchOutlined
                style={{
                    color: filtered ? '#1677ff' : undefined,
                }}
            />
        ),
        onFilter: (value, record) =>
            record[dataIndex].toString().toLowerCase().includes(value.toLowerCase()),
    });

    // 表格列配置
    const columns = [
        {
            title: '封禁IP',
            dataIndex: 'banIP',
            key: 'banIP',
            width: '10%',
            ...getColumnSearchProps('banIP'),
        },
        {
            title: '触发规则ID',
            dataIndex: 'banRuleID',
            key: 'banRuleID',
            width: '10%',
            ...getColumnSearchProps('banRuleID'),
        },
        {
            title: '封禁时长',
            dataIndex: 'banPeriod',
            key : 'banPeriod',
            width: '15%',
            ...getColumnSearchProps('banPeriod'),
        },
        {
            title: '解禁时间',
            dataIndex: 'unbanTime',
            key : 'unbanTime',
            ...getColumnSearchProps('unbanTime'),
        },
        {
            title: '描述',
            dataIndex: 'describe',
            key: 'describe',
            width: '40%',
            ...getColumnSearchProps('describe'),
        },
        {
            title: 'Action',
            dataIndex: 'operation',
            render: (_, record) => (
                <Popconfirm
                    title="Are you sure to delete this IP?"
                    onConfirm={() => handleDelete(record.banIP)}
                >
                    <a>Delete</a>
                </Popconfirm>
            ),
        },
    ];

    // 处理分页
    const handleTableChange = (pagination) => {
        setPagination(pagination);
    };

    // 删除操作
    const handleDelete = async (banIP) => {
        const response = await fetch(`http://localhost:8080/DeleteBanIP?banIP=${banIP}`, {
            method: 'GET',
        });
        if (response.ok) {
            await fetchData(searchText, pagination.current, pagination.pageSize);
        }
    };

    const handleAdd = async (values) => {
        // 将 banRuleID 转换为整数
        const newValues = {
            ...values,
            banRuleID: parseInt(values.banRuleID, 10), // 转换为整数
            unbanTime: parseInt(values.unbanTime, 10),
            banPeriod: parseInt(values.unbanTime,10)
        };

        // 确保转换成功（如果不是有效的数字，默认设置为 0）
        if (isNaN(newValues.banRuleID)) {
            newValues.banRuleID = 0;
        }
        // 确保转换成功（如果不是有效的数字，默认设置为 15）
        if (isNaN(newValues.unbanTime)) {
            newValues.unbanTime = 15;
        }
        // 确保转换成功（如果不是有效的数字，默认设置为 15）
        if (isNaN(newValues.banPeriod)) {
            newValues.unbanTime = 15;
        }

        // 将 banTime 转换为当前时间加指定分钟后的日期时间（ISO 格式）
        const now = new Date();
        now.setMinutes(now.getMinutes() + newValues.unbanTime);
        now.setHours(now.getHours()+8)
        newValues.unbanTime = now.toISOString();

        // 发送请求
        const response = await fetch('http://localhost:8080/AddBanIP', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(newValues),
        });

        if (response.ok) {
            await fetchData(searchText, pagination.current, pagination.pageSize);
            setOpenModal(false); // 关闭模态框
        }
    };

    // 打开模态框
    const showModal = () => {
        setOpenModal(true);
    };

    const handleCancel = () => {
        setOpenModal(false);
    };

    return (
        <>
            <Button
                type="primary"
                onClick={showModal}
                style={{ marginBottom: 16 }}
            >
                Add Ban IP
            </Button>
            <Table
                columns={columns}
                dataSource={data}
                rowKey="banIP"
                pagination={{
                    current: pagination.current,
                    pageSize: pagination.pageSize,
                    total,
                    showTotal: (total) => `Total ${total} items`,
                }}
                onChange={handleTableChange}
            />
            <Modal
                title="添加封禁IP"
                visible={openModal}
                onCancel={handleCancel}
                footer={null}
            >
                <Form
                    form={form}
                    onFinish={handleAdd}
                    initialValues={{ banIP: '', banRuleID: '', describe: '' }}
                >
                    <Form.Item
                        label="封禁IP"
                        name="banIP"
                        rules={[{ required: true, message: 'Please input the IP to ban!' }]}
                    >
                        <Input />
                    </Form.Item>
                    <Form.Item
                        label="描述"
                        name="describe"
                        rules={[{ required: true, message: 'Please input the description!' }]}
                    >
                        <Input />
                    </Form.Item>
                    <Form.Item
                        label="封禁时长"
                        name="unbanTime"
                        rules={[{ required: true, message: 'Please input the BanTime!' }]}
                    >
                        <Input />
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" htmlType="submit">
                            提交
                        </Button>
                    </Form.Item>
                </Form>
            </Modal>
        </>
    );
};

