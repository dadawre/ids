"use client";
import React, { useEffect, useState, useRef } from "react";
import { Table, Input, Space, Button, Pagination } from "antd";
import { SearchOutlined } from "@ant-design/icons";
import Highlighter from "react-highlight-words";
import axios from "axios";

export default function RecordTable() {
    const [data, setData] = useState([]);
    const [searchText, setSearchText] = useState("");
    const [searchedColumn, setSearchedColumn] = useState("");
    const searchInput = useRef(null);

    const [filters, setFilters] = useState({});
    const [pagination, setPagination] = useState({
        current: 1,
        pageSize: 10,
        total: 0,
    });

    const fetchData = () => {
        axios
            .get("http://localhost:8080/logData", {
                params: {
                    ...filters,
                    page: pagination.current,
                    pageSize: pagination.pageSize,
                },
            })
            .then((response) => {
                setData(response.data.items);
                setPagination((prevPagination) => ({
                    ...prevPagination,
                    total: response.data.total,
                }));
            })
            .catch((error) => {
                console.error("Error fetching data", error);
            });
    };

    useEffect(() => {
        fetchData();
    }, [filters, pagination.current, pagination.pageSize]);

    const handleSearch = (selectedKeys, confirm, dataIndex) => {
        confirm();
        setSearchText(selectedKeys[0]);
        setSearchedColumn(dataIndex);
        setFilters((prev) => ({ ...prev, [dataIndex]: selectedKeys[0] }));
    };

    const handleReset = (clearFilters, dataIndex) => {
        clearFilters();
        setSearchText("");
        setFilters((prev) => ({ ...prev, [dataIndex]: "" }));
    };

    const getColumnSearchProps = (dataIndex) => ({
        filterDropdown: ({ setSelectedKeys, selectedKeys, confirm, clearFilters, close }) => (
            <div style={{ padding: 8 }} onKeyDown={(e) => e.stopPropagation()}>
                <Input
                    ref={searchInput}
                    placeholder={`Search ${dataIndex}`}
                    value={selectedKeys[0]}
                    onChange={(e) => setSelectedKeys(e.target.value ? [e.target.value] : [])}
                    onPressEnter={() => handleSearch(selectedKeys, confirm, dataIndex)}
                    style={{ marginBottom: 8, display: "block" }}
                />
                <Space>
                    <Button
                        type="primary"
                        onClick={() => handleSearch(selectedKeys, confirm, dataIndex)}
                        icon={<SearchOutlined />}
                        size="small"
                        style={{ width: 90 }}
                    >
                        Search
                    </Button>
                    <Button
                        onClick={() => clearFilters && handleReset(clearFilters, dataIndex)}
                        size="small"
                        style={{ width: 90 }}
                    >
                        Reset
                    </Button>
                    <Button type="link" size="small" onClick={() => close()}>
                        Close
                    </Button>
                </Space>
            </div>
        ),
        filterIcon: (filtered) => <SearchOutlined style={{ color: filtered ? "#1677ff" : undefined }} />,
        onFilter: (value, record) => record[dataIndex]?.toString().toLowerCase().includes(value.toLowerCase()),
        render: (text) =>
            searchedColumn === dataIndex ? (
                <Highlighter
                    highlightStyle={{ backgroundColor: "#ffc069", padding: 0 }}
                    searchWords={[searchText]}
                    autoEscape
                    textToHighlight={text ? text.toString() : ""}
                />
            ) : (
                text
            ),
    });

    const columns = [
        { title: "ID", dataIndex: "ID", key: "ID",width: "7%", ...getColumnSearchProps("ID"), fixed: "top" },
        { title: "Src IP", dataIndex: "SrcIPAddr",width: "15%", key: "SrcIPAddr", ...getColumnSearchProps("SrcIPAddr") },
        { title: "Dst IP", dataIndex: "DstIPAddr",width: "15%", key: "DstIPAddr", ...getColumnSearchProps("DstIPAddr") },
        { title: "Src Port", dataIndex: "SrcPortAddr",width: "12%", key: "SrcPortAddr", ...getColumnSearchProps("SrcPortAddr") },
        { title: "Dst Port", dataIndex: "DstPortAddr",width: "12%", key: "DstPortAddr", ...getColumnSearchProps("DstPortAddr") },
        { title: "Protocol", dataIndex: "TopProtocol",width: "12%", key: "TopProtocol", ...getColumnSearchProps("TopProtocol") },
        { title: "Rule ID", dataIndex: "RuleID",width: "10%", key: "RuleID", ...getColumnSearchProps("RuleID") },
        { title: "Rule Name", dataIndex: "RuleName",width: "15%", key: "RuleName", ...getColumnSearchProps("RuleName") },
        { title: "Message", dataIndex: "DescribeMessage",width: "25%", key: "DescribeMessage", ...getColumnSearchProps("DescribeMessage") },
        { title: "Trigger Time", dataIndex: "TriggerTime",width: "15%", key: "TriggerTime", ...getColumnSearchProps("TriggerTime"), fixed: "top" },
        { title: "Action", dataIndex: "ActionType",width: "10%", key: "ActionType", ...getColumnSearchProps("ActionType") },
        { title: "TCP Flag", dataIndex: "TcpFlag",width: "10%", key: "TcpFlag", ...getColumnSearchProps("TcpFlag") },
        { title: "Packet Length", dataIndex: "PacketLength",width: "15%", key: "PacketLength", ...getColumnSearchProps("PacketLength"), fixed: "bottom" },
    ];

    const handlePageChange = (page, pageSize) => {
        setPagination({ current: page, pageSize, total: pagination.total });
    };

    return (
        <div>
            <Table
                columns={columns}
                dataSource={data}
                rowKey="ID"
                pagination={{
                    current: pagination.current,
                    pageSize: pagination.pageSize,
                    total: pagination.total,
                    onChange: handlePageChange,
                    showSizeChanger: true,
                    showQuickJumper: true,
                }}
                scroll={{ x: "125%", y: 500 }}
                sticky
            />
        </div>
    );
}
