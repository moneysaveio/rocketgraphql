import { Form, Input, Button, Descriptions } from 'antd';
import { Row } from "antd";
import  axios from "axios";
import { useState } from 'react';
import { useHistory } from 'react-router';
import { useAppDispatch, useAppSelector } from '../app/hooks';

const Signup = () => {
  const history = useHistory();
  const [user, setUser] = useState(null);
  const dispatch = useAppDispatch();
  const userFromStore = useAppSelector((state: any) => state.user.value)

  const onFinish = (values: any) => {
    // post to API here
    const API_URL = `${process.env.REACT_APP_API_ENDPOINT}/signup`;
    axios.post(API_URL, JSON.stringify(values), {
        headers: {
            // 'application/json' is the modern content-type for JSON, but some
            // older servers may use 'text/json'.
            // See: http://bit.ly/text-json
            'Content-Type': 'application/json',
        }
    }).then(response => {
      console.log("bare response:", response);
      dispatch({
        type: "user/set",
        payload: response.data.Email
      })
      console.log('Successfully added user to state:', userFromStore);
      setUser(userFromStore);
      localStorage.setItem('user', JSON.stringify(response.data));
      history.push("/dashboard");
    })
    console.log('Success:', values, API_URL);
  };

  const onFinishFailed = (errorInfo: any) => {
    console.log('Failed:', errorInfo);
  };

  return (
    <Row justify="center" align="middle" style={{minHeight: '100vh'}}>
        <Form
        id="login-form"
        name="basic"
        labelCol={{ span: 8 }}
        wrapperCol={{ span: 16 }}
        initialValues={{ remember: true }}
        onFinish={onFinish}
        onFinishFailed={onFinishFailed}
        >
        <Form.Item
            label="Email"
            name="email"
            rules={[{ required: true, message: 'Please input your username!' }]}
        >
            <Input />
        </Form.Item>

        <Form.Item
            label="Password"
            name="password"
            rules={[{ required: true, message: 'Please input your password!' }]}
        >
            <Input.Password />
        </Form.Item>

        <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
            <Button type="primary" htmlType="submit" id="register-button">
            Submit
            </Button>
        </Form.Item>
        </Form>
    </Row>
  );
};

export default Signup;

function userState(arg0: null): [any, any] {
  throw new Error('Function not implemented.');
}
