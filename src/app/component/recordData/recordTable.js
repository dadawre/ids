"use client";
import React, { useEffect, useState, useRef } from 'react';
import { Table, Input, Space, Button, Pagination } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import Highlighter from 'react-highlight-words';
import axios from 'axios';

export default function RecordTable() {
    const [data, setData] = useState([]);
    const [searchText, setSearchText] = useState('');
    const [searchedColumn, setSearchedColumn] = useState('');
    const searchInput = useRef(null);

    // 用于保存搜索条件
    const [filters, setFilters] = useState({
        DescribeMessage: '',
        SrcAddr: '',
        DstAddr: '',
        ActionType: '',
    });

    // 用于分页
    const [pagination, setPagination] = useState({
        current: 1,
        pageSize: 10,
        total: 0,
    });

    // Fetch data from the backend with filters and pagination
    const fetchData = () => {
        axios.get('http://localhost:8080/logData', {
            params: {
                ...filters,
                page: pagination.current,
                pageSize: pagination.pageSize,
            },
        })
            .then(response => {
                setData(response.data.items); // Ensure this matches the structure of the response
                setPagination(prevPagination => ({
                    ...prevPagination,
                    total: response.data.total,  // Update total with the response from the backend
                }));
            })
            .catch(error => {
                console.error("There was an error fetching the data!", error);
            });
    };

    useEffect(() => {
        fetchData();
    }, [filters, pagination.current, pagination.pageSize]); // Re-fetch when filters or pagination change

    const handleSearch = (selectedKeys, confirm, dataIndex) => {
        confirm();
        setSearchText(selectedKeys[0]);
        setSearchedColumn(dataIndex);

        // Update filter conditions
        setFilters(prevFilters => ({
            ...prevFilters,
            [dataIndex]: selectedKeys[0],
        }));
    };

    const handleReset = (clearFilters, dataIndex) => {
        clearFilters();
        setSearchText('');
        setFilters(prevFilters => ({
            ...prevFilters,
            [dataIndex]: '',
        }));
    };

    const getColumnSearchProps = (dataIndex) => ({
        filterDropdown: ({ setSelectedKeys, selectedKeys, confirm, clearFilters, close }) => (
            <div
                style={{
                    padding: 8,
                }}
                onKeyDown={(e) => e.stopPropagation()}
            >
                <Input
                    ref={searchInput}
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
                        onClick={() => clearFilters && handleReset(clearFilters, dataIndex)}
                        size="small"
                        style={{
                            width: 90,
                        }}
                    >
                        Reset
                    </Button>
                    <Button
                        type="link"
                        size="small"
                        onClick={() => {
                            confirm({
                                closeDropdown: false,
                            });
                            setSearchText(selectedKeys[0]);
                            setSearchedColumn(dataIndex);
                        }}
                    >
                        Filter
                    </Button>
                    <Button
                        type="link"
                        size="small"
                        onClick={() => {
                            close();
                        }}
                    >
                        Close
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
        filterDropdownProps: {
            onOpenChange(open) {
                if (open) {
                    setTimeout(() => searchInput.current?.select(), 100);
                }
            },
        },
        render: (text) =>
            searchedColumn === dataIndex ? (
                <Highlighter
                    highlightStyle={{
                        backgroundColor: '#ffc069',
                        padding: 0,
                    }}
                    searchWords={[searchText]}
                    autoEscape
                    textToHighlight={text ? text.toString() : ''}
                />
            ) : (
                text
            ),
    });

    const columns = [
        {
            title: 'Describe Message',
            dataIndex: 'DescribeMessage',
            key: 'DescribeMessage',
            width: '30%',
            ...getColumnSearchProps('DescribeMessage'),
        },
        {
            title: 'Action Type',
            dataIndex: 'ActionType',
            key: 'ActionType',
            width: '20%',
            ...getColumnSearchProps('ActionType'),
        },
        {
            title: 'Source Address',
            dataIndex: 'SrcAddr',
            key: 'SrcAddr',
            ...getColumnSearchProps('SrcAddr'),
        },
        {
            title: 'Destination Address',
            dataIndex: 'DstAddr',
            key: 'DstAddr',
            ...getColumnSearchProps('DstAddr'),
        },
    ];

    const handlePageChange = (page, pageSize) => {
        setPagination({
            current: page,
            pageSize: pageSize,
            total: pagination.total, // Preserve total
        });
    };

    return (
        <div>
            <Table
                columns={columns}
                dataSource={data}
                rowKey="key"
                pagination={{
                    current: pagination.current,
                    pageSize: pagination.pageSize,
                    total: pagination.total,
                    onChange: handlePageChange,
                    showSizeChanger: true,
                    showQuickJumper: true,
                }} // Use Ant Design pagination
            />
        </div>
    );
}
